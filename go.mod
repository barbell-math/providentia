module code.barbellmath.net/barbell-math/providentia

go 1.25.1

require (
	code.barbellmath.net/barbell-math/smoothbrain-argparse v0.0.0-20251025080248-f6ed35565943
	code.barbellmath.net/barbell-math/smoothbrain-bs v0.0.0-20250831070729-5632f2448d7f
	code.barbellmath.net/barbell-math/smoothbrain-cgoGlue v0.0.0-20250901064600-3ef6fc548d82
	code.barbellmath.net/barbell-math/smoothbrain-cgoTest v0.0.0-20251112072621-434fffc6deb0
	code.barbellmath.net/barbell-math/smoothbrain-csv v0.0.0-20251025073939-2b81a8dfdeac
	code.barbellmath.net/barbell-math/smoothbrain-errs v0.0.0-20251025073532-bd7907d34fd8
	code.barbellmath.net/barbell-math/smoothbrain-jobQueue v0.0.0-20251025073702-62476ee76303
	code.barbellmath.net/barbell-math/smoothbrain-logging v0.0.0-20250831071438-e5e955330baf
	code.barbellmath.net/barbell-math/smoothbrain-sqlmigrate v0.0.0-20250831071515-53ff01a18bae
	code.barbellmath.net/barbell-math/smoothbrain-test v0.0.0-20250831071138-3f0f71428ad0
	github.com/jackc/pgx/v5 v5.7.6
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/exp v0.0.0-20251125195548-87e1e737ad39 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)

replace code.barbellmath.net/barbell-math/smoothbrain-cgoTest => ../smoothbrain-cgoTest // TODO - remove!
