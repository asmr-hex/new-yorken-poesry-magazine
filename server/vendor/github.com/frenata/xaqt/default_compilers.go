package xaqt

// the compilers included with the xaqt package.
// TODO (cw|8.26.2018) auto marshall into json when tests are run...this way
// users who want to use their own compilers.json will have a useful template.
var DEFAULT_COMPILERS = Compilers{
	"python": CompilerDetails{
		ExecutionDetails{
			Compiler:   "python",
			SourceFile: "file.py",
		},
		CompositionDetails{
			Boilerplate:   "import sys\ninput = sys.stdin.read()\nargs = input.split('\\n')\n",
			CommentPrefix: "#",
		},
	},
	"golang": CompilerDetails{
		ExecutionDetails{
			Compiler:   "go run",
			SourceFile: "file.go",
		},
		CompositionDetails{
			Boilerplate:   "package main\nimport (\n\t\"os\"\n\t\"bufio\"\n\t\"strings\"\n\t\"fmt\"\n)\nvar input, _ = bufio.NewReader(os.Stdin).ReadString('\\a')\nvar args = strings.Split(input, \"\\n\")\n\nfunc main(){}",
			CommentPrefix: "//",
		},
	},
	"ruby": CompilerDetails{
		ExecutionDetails{
			Compiler:   "ruby",
			SourceFile: "file.rb",
		},
		CompositionDetails{
			CommentPrefix: "#",
		},
	},
	"clojure": CompilerDetails{
		ExecutionDetails{
			Compiler:   "clojure",
			SourceFile: "file.clj",
		},
		CompositionDetails{
			CommentPrefix: ";",
		},
	},
	"javascript": CompilerDetails{
		ExecutionDetails{
			Compiler:   "nodejs",
			SourceFile: "file.js",
		},
		CompositionDetails{
			CommentPrefix: "//",
		},
	},
	"c++": CompilerDetails{
		ExecutionDetails{
			Compiler:           "g++ -o /usercode/a.out",
			SourceFile:         "file.cpp",
			OptionalExecutable: "/usercode/a.out",
		},
		CompositionDetails{
			CommentPrefix: "//",
		},
	},
	"java": CompilerDetails{
		ExecutionDetails{
			Compiler:           "javac",
			SourceFile:         "file.java",
			OptionalExecutable: "/entrypoint/javaRunner.sh",
		},
		CompositionDetails{
			CommentPrefix: "//",
		},
	},
	"perl": CompilerDetails{
		ExecutionDetails{
			Compiler:   "perl",
			SourceFile: "file.pl",
		},
		CompositionDetails{
			CommentPrefix: "#",
		},
	},
	"c#": CompilerDetails{
		ExecutionDetails{
			Compiler:           "mcs",
			SourceFile:         "file.cs",
			OptionalExecutable: "mono /usercode/file.exe",
		},
		CompositionDetails{
			CommentPrefix: "//",
		},
	},
	"bash": CompilerDetails{
		ExecutionDetails{
			Compiler:   "/bin/bash",
			SourceFile: "file.sh",
		},
		CompositionDetails{
			CommentPrefix: "#",
		},
	},
	"haskell": CompilerDetails{
		ExecutionDetails{
			Compiler:           "ghc -o /usercode/a.out",
			SourceFile:         "file.hs",
			OptionalExecutable: "/usercode/a.out",
		},
		CompositionDetails{
			CommentPrefix: "--",
		},
	},
	"scala": CompilerDetails{
		ExecutionDetails{
			Compiler:   "scala",
			SourceFile: "file.scala",
			Disabled:   true,
		},
		CompositionDetails{},
	},
	"php": CompilerDetails{
		ExecutionDetails{
			Compiler:   "php",
			SourceFile: "file.php",
		},
		CompositionDetails{},
	},
	"rust": CompilerDetails{
		ExecutionDetails{
			Compiler:           "'/opt/rust/.cargo/bin/rustc'",
			SourceFile:         "file.rs",
			OptionalExecutable: "/usercode/a.out",
			CompilerFlags:      "'-o /usercode/a.out'",
			Disabled:           true,
		},
		CompositionDetails{},
	},
}
