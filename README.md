# Github data processor with TopK output

## A readme on how to run the solution ##
How to download: 

- Checkout/download this set of files to any folder
- This repo already contains `data` folder with input files - keep it to avoid unnecessary actions

### How to run
```shell
make build
```
After that binary at path `./bin/app` will be created.
Following commands depend on this binary.

#### Top 10 active users sorted by amount of PRs created and commits pushed
This is default scenario.
It can be executed  by following:
```shell
make top1
  >  Top 10 active users sorted by amount of PRs created and commits pushed...
LombiqBot                 1529
renovate[bot]             535
pull[bot]                 384
direwolf-github           341
lihkg-boy                 331
ripamf2991                311
renovate-bot              232
otiny                     222
dependabot[bot]           183
dependabot-preview[bot]   155
```

#### Top 10 repositories sorted by amount of commits pushed
```shell
make top2
  >  Top 10 repositories sorted by amount of commits pushed...
lihkg-backup/thread                      331
otiny/up                                 222
ripamf2991/ntdtv                         167
ripamf2991/djy                           139
wessilfie/wessilfie.github.io            108
Lombiq/Orchard                           96
himobi/hotspot                           90
wigforss/java-8-base                     87
geos4s/18w856162                         79
SmartThingsCommunity/SmartThingsPublic   68
```

#### Top 10 repositories sorted by amount of watch events
```shell
make top3
  >  Top 10 repositories sorted by amount of watch events...
victorqribeiro/isocity                44
GitHubDaily/GitHubDaily               11
neutraltone/awesome-stock-resources   11
sw-yx/spark-joy                       10
imsnif/bandwhich                      8
Chakazul/Lenia                        7
BurntSushi/xsv                        7
neeru1207/AI_Sudoku                   6
ErikCH/DevYouTubeList                 6
testerSunshine/12306                  6
```

To use custom parameters, run ./bin/app.
The binary is created with cobra CLI, so help is available:
```shell
Usage:
  app top [flags]

Flags:
      --entity_entity_column_index int       
      --entity_file string                    (default "./data/actors.csv")
      --entity_name_column_index int          (default 1)
      --event_types strings                   (default [PushEvent,PullRequestEvent])
      --events_entity_column_index int        (default 2)
      --events_event_type_column_index int    (default 1)
      --events_file string                    (default "./data/events.csv")
  -h, --help                                 help for top
      --k uint32                              (default 10)
```

### How it works
- For TopK task https://github.com/migotom/heavykeeper is used in single-worker mode, otherwise the results are unstable. 
Heavy Keeper algorithm seems to be better than Sorted Set - https://redis.com/blog/meet-top-k-awesome-probabilistic-addition-redisbloom/ .
- CLI interface provides all necessary parameters to customize processing (all csv file indexes, source file names, etc)

### Structure
- main.go that executes root command of cobra.
- `./command/top.go` - `top` command which is root subcommand responsible for all top functionality.
Top command is responsible to collect flag values and call CsvParserApp.
`CsvParserAppBuilder` is needed to replace build logic in tests.
- `./app/csv.go` - processing logic. Consists of 3 function calls described below.
### Processing logic
- `getIDsTop` initializes HeavyKeeper instance, reads events file and on calls HeavyKeeper on every matching line.
Returns top on K (10 by default) leader pairs of entity ID + count.
- `getNamedTop` uses leaders list from previous step to prepare named list.
Names are being read from entity file (actors.csv or repos.csv). All top items are iterated for every entity file line until
all names are found or EOF.
- `WriteResults` writes leaderboard list to `io.Writer` implementation passed to `CsvParserApp` on build stage.