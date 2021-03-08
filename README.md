MAD
===

MaD is a markdown preprocessor, it aims to create the best writing enviroment for every use case.

MaD works with little scripts that extendends the Markdown language
```
[math]: #(2x + 2)
```
For example calls the `math` script with the argument `2x + 2`

The location of all scripts in by default `$HOME/.config/mad/bin/` but can be changed by `MAD_PATH`

The script is called with the first argument being all the given argument and then the same argument splitted by spaces.
```
$HOME/.config/mad/bin/math "2x + 2" 2x + 2
```

### Roadmap

- [ ] Multi-line input
- [ ] Preview
- [ ] Multiple Output
- [ ] Better syntax
- [ ] Better parser
