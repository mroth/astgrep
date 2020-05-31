# astgrep

Quick AST-aware tool to run regex pattern searches against only the contents of
all* string literals within Go source code files. This is useful if you're
trying to find all matches of a pattern that also commonly appears in your
source code, and want to avoid the noise.

(*: By default, string literals that are part of import statements are ignored.)

