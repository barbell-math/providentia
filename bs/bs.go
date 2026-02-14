package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	sbbs "code.barbellmath.net/barbell-math/smoothbrain-bs"
)

var (
	ffmpegConf = []string{
		// Don't need ffplay or it's deps
		"--disable-ffplay", "--disable-libxcb", "--disable-libxcb-shm", "--disable-libxcb-xfixes", "--disable-libxcb-shape", "--disable-sdl2",
		// Re-enable VAAPI for hw accel on linux if desired, but know that
		// enabling VAAPI will create wayland and x11 dependencies.
		"--disable-vaapi", "--disable-xlib",
		// Enable things that are under a gpl license as well as vulkan
		"--enable-gpl", "--enable-vulkan",
		`--extra-cflags="-I${VULKAN_SDK}/include"`,
		`--extra-ldflags="-L${VULKAN_SDK}/lib"`,
		// Food for thought: might be useful one day for advanced vulkan filtering
		// "--enable-libplacebo",
	}
)

func main() {
	setupDevTarget()
	setupDownloadTargets()
	setupBuildFFmpegTarget()
	sbbs.RegisterInstallBashAutocompleteTarget()
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterDepTargets()
	sbbs.RegisterGoMarkDocTargets(sbbs.NewReadmeOpts().
		SetRepoUrl("dummy").
		SetDirsToRead("./logic").
		SetPostStages(sbbs.SedStage("dummy/blob/main", "", "README.md")),
	)
	sbbs.RegisterGoEnumTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.AllGoTargets().
		SetEnvVars(map[string]string{"CC": "clang-21"}),
	)
	sbbs.RegisterMergegateTarget(sbbs.NewMergegateTargets().
		SetPreStages(sbbs.TargetAsStage("install.goenum")),
	)
	sbbs.Main()
}

func setupDevTarget() {
	sbbs.RegisterTarget(
		context.Background(),
		"configure.dev",
		sbbs.CdToRepoRoot(),
		sbbs.TargetAsStage("download.submodules"),
		sbbs.TargetAsStage("download.vulkan"),
		sbbs.TargetAsStage("build.ffmpeg"),
	)
}

func setupDownloadTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"download.submodules",
		sbbs.CdToRepoRoot(),
		sbbs.Stage(
			"Clone submodules",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := sbbs.RunStdout(
					ctxt, "git", "submodule", "init",
				); err != nil {
					return err
				}
				return sbbs.RunStdout(
					ctxt, "git", "submodule", "update", "--depth", "1",
				)
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"download.vksdk",
		sbbs.Stage(
			"Cd to repo root",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				repoRoot, err := sbbs.GitRevParse(ctxt)
				if err != nil {
					return err
				}
				return sbbs.Cd(repoRoot)
			},
		),
		sbbs.Stage(
			"cd to vkSdk src dir",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.Cd(path.Join("_deps", "vkSdk"))
			},
		),
		sbbs.Stage(
			"Clear vkSdk dir",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RmDir("./vulkan-linux-x86_64-1.4.341.1")
			},
		),
		sbbs.Stage(
			"Download sdk",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(
					ctxt, "wget",
					"https://sdk.lunarg.com/sdk/download/1.4.341.1/linux/vulkansdk-linux-x86_64-1.4.341.1.tar.xz",
				)
			},
		),
		sbbs.Stage(
			"Extract sdk",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(
					ctxt, "tar", "xf", "vulkansdk-linux-x86_64-1.4.341.1.tar.xz",
				)
			},
		),
		sbbs.Stage(
			"Remove sdk tar",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(
					ctxt, "rm", "vulkansdk-linux-x86_64-1.4.341.1.tar.xz",
				)
			},
		),
	)
}

func setupBuildFFmpegTarget() {
	repoRoot := ""

	sbbs.RegisterTarget(
		context.Background(),
		"build.ffmpeg",
		sbbs.Stage(
			"Cd to repo root",
			func(ctxt context.Context, cmdLineArgs ...string) (err error) {
				repoRoot, err = sbbs.GitRevParse(ctxt)
				if err != nil {
					return
				}
				return sbbs.Cd(repoRoot)
			},
		),
		sbbs.Stage(
			"cd to ffmpeg src dir",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.Cd(path.Join("_deps", "ffmpeg", "src"))
			},
		),
		sbbs.Stage(
			"Configure ffmpeg",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				buildLogPath := path.Join(repoRoot, "bs", "logs", "ffmpegBuild.log")
				f, err := os.Create(buildLogPath)
				if err != nil {
					return err
				}
				defer f.Close()

				clangDir, err := exec.LookPath("clang")
				if err != nil {
					return err
				}

				configureScript := fmt.Sprintf(
					`source %s
					./configure --prefix=%s --cc=%s %s`,
					path.Join(repoRoot, "_deps", "vkSdk", "1.4.341.1", "setup-env.sh"),
					path.Join(repoRoot, "_deps", "ffmpeg"),
					clangDir,
					strings.Join(ffmpegConf, " "),
				)
				return sbbs.RunBashScript(ctxt, f, "", configureScript)
			},
		),
		sbbs.Stage(
			"Make ffmpeg",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(
					ctxt, "make", fmt.Sprintf("-j%d", runtime.NumCPU()),
				)
			},
		),
		sbbs.Stage(
			"Install ffmpeg",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "make", "install")
			},
		),
	)
}
