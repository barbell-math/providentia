package prov_exportedbs

import (
	"fmt"
	"path"
	"runtime"
	"strconv"

	"code.barbellmath.net/carmichaeljr/smoothbrain/sbbs"
)

type (
	Conf struct {
		FFmpegConf        []string
		RequiredVkSdkLibs []string
		VulkanSdkVersion  string
		VulkanEnvVars     []sbbs.EnvVar
	}
)

var (
	moduleRoot = path.Join(sbbs.CallerDir(), "..", "..")

	FmtStages        = sbbs.DefaultGoFmtStages[Conf](moduleRoot)
	UpdateDepsStages = sbbs.DefaultGoUpdatePkgsStages[Conf](moduleRoot)
	GenStages        = []sbbs.StageFunc[Conf]{
		sbbs.TmpSetEnvVarStage(
			[]sbbs.EnvVar{{Name: "CC", Value: "clang"}},
			sbbs.DefaultRunStage[Conf](moduleRoot, "go", "gen", "./..."),
		),
	}
	TestStages = []sbbs.StageFunc[Conf]{
		sbbs.TmpSetEnvVarStage(
			[]sbbs.EnvVar{{Name: "CC", Value: "clang"}},
			sbbs.DefaultRunStage[Conf](moduleRoot, "go", "test", "-v", "./..."),
		),
	}

	MergegateStages = sbbs.DefaultMergegate(sbbs.Mergegate[Conf]{
		ModuleRoot: moduleRoot,
		FmtStages:  FmtStages,
		GenStages:  GenStages,
		TestStages: TestStages,
	})
	DevSetup = []sbbs.StageFunc[Conf]{
		SubModuleSetup,
		VkSdkSetup,
		FFmpegSetup,
	}
)

func NewConf() Conf {
	return Conf{
		VulkanSdkVersion: "1.4.341.1",
		FFmpegConf: []string{
			// Don't need ffplay or it's deps
			"--disable-ffplay",
			"--disable-libxcb",
			"--disable-libxcb-shm",
			"--disable-libxcb-xfixes",
			"--disable-libxcb-shape",
			"--disable-sdl2",
			// Re-enable VAAPI for hw accel on linux if desired, but know that
			// enabling VAAPI will create wayland and x11 dependencies.
			"--disable-vaapi",
			"--disable-xlib",
			// Enable things that are under a gpl license as well as vulkan
			"--enable-gpl",
			"--enable-vulkan",
			"--enable-libglslang",
			// "--enable-libx264",
			`--extra-cflags="-I${VULKAN_SDK}/include"`,
			`--extra-ldflags="-L${VULKAN_SDK}/lib"`,
		},
		RequiredVkSdkLibs: []string{
			"libglslang.a",
			"libglslang-default-resource-limits.a",
			"libOSDependent.a",
			"libMachineIndependent.a",
			"libGenericCodeGen.a",
			"libSPIRV.a",
			"libSPIRV-Tools.a",
			"libshaderc_combined.a",
		},
	}
}

func SubModuleSetup(_ *Conf, output sbbs.Output) error {
	if err := sbbs.RunWithOpts(&sbbs.RunOpts{
		Cwd:    moduleRoot,
		Output: output,
	}, "git", "submodule", "init"); err != nil {
		return err
	}
	if err := sbbs.RunWithOpts(&sbbs.RunOpts{
		Cwd:    moduleRoot,
		Output: output,
	}, "git", "submodule", "update", "--depth", "1"); err != nil {
		return err
	}
	return nil
}

func VulkanEnvVars(stages ...sbbs.StageFunc[Conf]) sbbs.StageFunc[Conf] {
	var resetTmps func()
	var resetPaths func()
	rv := []sbbs.StageFunc[Conf]{func(conf *Conf, output sbbs.Output) (err error) {
		vkSdkPath := path.Join(
			moduleRoot, "_deps", "vkSdk", conf.VulkanSdkVersion, "x86_64",
		)
		if resetTmps, err = sbbs.TmpSetEnvVar(output, sbbs.EnvVar{
			Name:  "VULKAN_SDK",
			Value: vkSdkPath,
		}, sbbs.EnvVar{
			Name:  "VK_LAYER_PATH",
			Value: path.Join(vkSdkPath, "share", "vulkan", "explicit_layer.d"),
		}, sbbs.EnvVar{
			Name:  "VK_ADD_LAYER_PATH",
			Value: path.Join(vkSdkPath, "share", "vulkan", "explicit_layer.d"),
		}); err != nil {
			return
		}
		if resetPaths, err = sbbs.TmpAddPathToEnvVar(output, sbbs.PathListEnvVar{
			Name:  "PATH",
			Paths: []string{path.Join(vkSdkPath, "bin")},
		}, sbbs.PathListEnvVar{
			Name:  "PKG_CONFIG_PATH",
			Paths: []string{path.Join(vkSdkPath, "lib", "pkgconfig")},
		}, sbbs.PathListEnvVar{
			Name:  "LD_LIBRARY_PATH",
			Paths: []string{path.Join(vkSdkPath, "lib")},
		}); err != nil {
			return
		}
		return
	}}
	rv = append(rv, stages...)
	rv = append(rv, func(conf *Conf, output sbbs.Output) error {
		if resetTmps != nil {
			resetTmps()
		}
		if resetPaths != nil {
			resetPaths()
		}
		return nil
	})
	return sbbs.CoalesceStages("Vulkan Env Vars", rv...)
}

func VkSdkSetup(conf *Conf, output sbbs.Output) error {
	vkSdkPath := path.Join(moduleRoot, "_deps", "vkSdk")
	vkLibPath := path.Join(vkSdkPath, "lib")
	vkSdkVersionPath := path.Join(vkSdkPath, conf.VulkanSdkVersion)
	vkSdkTarPath := vkSdkVersionPath + ".tar.xz"
	vkSdkUrl := fmt.Sprintf(
		"https://sdk.lunarg.com/sdk/download/%s/linux/vulkansdk-linux-x86_64-%s.tar.xz",
		conf.VulkanSdkVersion, conf.VulkanSdkVersion,
	)

	if err := sbbs.Rm(output, vkSdkVersionPath); err != nil {
		return err
	}
	if err := sbbs.DownloadFile(output, vkSdkUrl, vkSdkTarPath); err != nil {
		return err
	}
	if err := sbbs.RunWithOpts(&sbbs.RunOpts{
		Output: output,
		Cwd:    vkSdkPath,
	}, "tar", "-xf", vkSdkTarPath); err != nil {
		return err
	}
	if err := sbbs.Rm(output, vkSdkTarPath); err != nil {
		return err
	}

	if err := sbbs.Rm(output, vkLibPath); err != nil {
		return err
	}
	if err := sbbs.MkDir(output, vkLibPath); err != nil {
		return err
	}
	for _, lib := range conf.RequiredVkSdkLibs {
		oldLibPath := path.Join(vkSdkVersionPath, "x86_64", "lib", lib)
		newLibPath := path.Join(vkLibPath, lib)
		if err := sbbs.Cp(output, oldLibPath, newLibPath); err != nil {
			return err
		}
	}
	return nil
}

var FFmpegSetup = VulkanEnvVars(func(conf *Conf, output sbbs.Output) error {
	ffmpegPath := path.Join(moduleRoot, "_deps", "ffmpeg")
	ffmpegSrcPath := path.Join(ffmpegPath, "src")

	buildLogPath := path.Join(moduleRoot, "bs", "logs", "ffmpegBuild.log")
	f, err := sbbs.Create(output, buildLogPath)
	if err != nil {
		return err
	}
	defer f.Close()

	args := []string{"./configure", "--prefix=" + ffmpegPath, "--cc=clang"}
	args = append(args, conf.FFmpegConf...)
	if err := sbbs.RunWithOpts(&sbbs.RunOpts{
		Output: output,
		Cwd:    ffmpegSrcPath,
	}, "bash", args...); err != nil {
		return err
	}

	if err := sbbs.RunWithOpts(&sbbs.RunOpts{
		Output: output,
		Cwd:    ffmpegSrcPath,
	}, "make", "-j", strconv.Itoa(runtime.NumCPU())); err != nil {
		return err
	}

	if err := sbbs.RunWithOpts(&sbbs.RunOpts{
		Output: output,
		Cwd:    ffmpegSrcPath,
	}, "make", "install"); err != nil {
		return err
	}

	return nil
})
