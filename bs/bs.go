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
		// For removing X11 depedency
		"--disable-ffplay", "--disable-libxcb", "--disable-libxcb-shm", "--disable-libxcb-xfixes", "--disable-libxcb-shape", "--disable-sdl2",
		// Enable things that are under a gpl license
		"--enable-gpl",
		// Enable vaapi for hw accel on linux
		"--enable-vaapi",
	}
)

func main() {
	setupDevTarget()
	setupBuildFFmpegTarget()
	sbbs.RegisterInstallBashAutocompleteTarget()
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterDepTargets()
	sbbs.RegisterGoMarkDocTargets(sbbs.NewReadmeOpts().
		SetRepoUrl("dummy").
		SetDirsToRead("./logic").
		SetPostStages(sbbs.SedStage("dummy/blob/main", "", "README.md")),
	)
	sbbs.RegisterSqlcTargets()
	sbbs.RegisterGoEnumTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.AllGoTargets().
		SetEnvVars(map[string]string{
			"CC": "clang-21",
		}),
	)
	sbbs.RegisterMergegateTarget(sbbs.NewMergegateTargets().
		SetPreStages(
			sbbs.TargetAsStage("install.goenum"),
			sbbs.TargetAsStage("install.sqlc"),
		),
	)
	sbbs.Main()
}

func setupDevTarget() {
	sbbs.RegisterTarget(
		context.Background(),
		"configure.dev",
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
					`./configure --prefix=%s --cc=%s %s`,
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
