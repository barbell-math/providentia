package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	sbbs "code.barbellmath.net/barbell-math/smoothbrain-bs"
)

var (
	ffmpegConf = []string{
		"--enable-gpl", "--enable-nonfree", "--enable-vaapi",
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

				configureScript := fmt.Sprintf(
					`./configure --prefix=%s %s`,
					path.Join(repoRoot, "_deps", "ffmpeg"),
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
