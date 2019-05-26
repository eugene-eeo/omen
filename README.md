# rtpreview

experimental 'real-time-previewer' for the terminal.
currently supports unicode well enough, but if the output of the command contains
colours it will not be handled properly.
core is complete for the most part.

    $ rtpreview -cmd 'ag "%s"'

todo:

 - [ ] testing
 - [ ] custom prompt
 - [ ] support for padding
 - [ ] custom colourscheme
 - [ ] handle colours
