# SRI

SRI is a cli tool to create sub ressource integrity hashes of a resource file.

## Usage

`sri [OPTION] <url>`

#### Options

- `-sha256` Use sha256 as hash function
- `-sha384` Use sha384 as hash function (default)
- `-sha512` Use sha512 as hash function

#### Examples

````
$ sri https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js

Returns:
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js" integrity="sha384-tsQFqpEReu7ZLhBV2VZlAu7zcOV+rXbYlF2cqB8txI/8aZajjp4Bqd+V6D5IgvKT" crossorigin="anonymous"></script>
````

````
$ sri -sha512 https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css

Returns:
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha512-tDXPcamuZsWWd6OsKFyH6nAqh/MjZ/5Yk88T5o+aMfygqNFPan1pLyPFAndRzmOWHKT+jSDzWpJv8krj6x1LMA==" crossorigin="anonymous">
````

## Installation

See [tags](https://git.bn4t.me/root/sri/tags) for available stable binaries. 

Alternatively all CI jobs store the latest binaries as artifacts. Artifacts are stored for 4 weeks.