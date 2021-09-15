# Prolific-Importer

A little tool for importing one or more stories into Pivotal Tracker via it's API.

Prolific-Importer takes story output from [Prolific](https://github.com/onsi/prolific) and imports it into [Pivotal Tracker](https://www.pivotaltracker.com/) without any manual steps.

## Installation

To install from source, make sure you have the Go 1.16+ toolchain installed, then:
`go install github.com/sneal/prolific-importer`

Or just download the OS X binary from the GitHub releases page.

## Usage

```bash
prolific-importer "$API_TOKEN" 12345 path/to/stories.csv
```

Will import the CSV as separate Tracker stories into the project with id 12345. You can find your API Token on your Tracker profile page and the project ID in your browser address bar or Tracker's **Projects Settings** page.

Prolific-Importer will also read content from standard input, which can be useful when combined with Prolific. For example:

```bash
prolific path/to/stories.prolific | prolific-importer "$API_TOKEN" 12345
```

## CSV Syntax

Stories are separated by newlines `\n` with the following fields:
- Title
- Type
- Description
- Labels

Here's an [example CSV](./fixtures/stories1.csv)
