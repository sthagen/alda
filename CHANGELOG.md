# CHANGELOG

## 2.0.1 (2021-07-05)

* Alda will now attempt to detect if it's running in an environment (e.g.
  the CMD program that ships with Windows 7) that does not support ANSI escape
  codes to display colored text. If the environment does not appear to support
  ANSI escape codes, Alda will not display colored text (which is better in that
  case because otherwise you would see a bunch of weird-looking characters in
  places where there should be colored text!).

* Prior to this release, it wasn't obvious that it's incorrect to enter a
  command like:

  ```
  alda play my-score.alda
  ```

  The correct way to specify a score file to read is to use the `-f, --file`
  option:

  ```
  alda play -f my-score.alda
  ```

  Instead of silently ignoring the provided file name, the Alda CLI will now
  print a useful error message.

## 2.0.0 (2021-06-30)

Alda 2 is a from-the-ground-up rewrite, optimized for simpler architecture,
better performance, and a foundation for future work to enable fun live coding
features.

For information about what's new, what's changed, and what to expect, check out
the [Alda 2 migration guide][migration-guide]!

[migration-guide]: https://github.com/alda-lang/alda/blob/master/doc/alda-2-migration-guide.md

---

## Earlier Versions

* [1.0.0 - 1.X.X](CHANGELOG-1.X.X.md)
* [0.1.0 - 0.X.X](CHANGELOG-0.X.X.md)
