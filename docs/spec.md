## What do I want to achieve

Each year I enjoy Advent Of Code. As an IC-turned-manager, I want to keep my skills sharp, _and_ 
I just enjoy it. But each year there are the same non-value-adding steps I take:

 - load the day's task in the browser
 - set up a new directory for that day
 - save the given test input into a text file
 - write boilerplate code in `solution.go`, incl.  reading input and passing it into a `solve` function
 - set up boilerplate test function against `solve` function
 - _actually_ get started solving. 


5 out of the 6 steps are repeated, can easily be automated and reduce time-to-value on 
solving the actual need. So, let's automate them. 

## Considerations

 - RATE LIMIT - Eric is a legend for his time and effort setting up AoC every year. Respect his infra and costs. <3 
 - Each invocation depends on `year`, `day`, a directory, but during a year, the year is constant.
 - Once you have X days reflected locally, the tool should know and instead only grab those available from X+1 onwards
 - There are only so many days available within a given year to be solved (1 released each day).
 - Solve for myself first, then package into a CLI tool for easier distribution.


## API

Some examples of possible evocations of the tool to help design interface.

- `aoc fetch` -- this should grab the configured (?) years days, based on what's already available on the fs
- `aoc fetch --year 2022 --day 4` -- this should specifically fetch day 4 of 2022 and save to fs
- `aoc fetch --year 2022 --day 4 --path ~/aoc-2022` -- this should do the same as above, but specific to that given path


### FS layout

- Allow user to determine the root dir.
- Within that root dir, each day fetched should get a new dir with left-padding on single-digit days for ordering (e.g. `day-04`)
- Use a config file (maybe `.aoc` in the user-configured root dir) that can contain the year, and signify the directory.

