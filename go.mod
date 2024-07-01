module github.com/dusmcd/pokedexcli

go 1.22.3

replace github.com/dusmcd/pokedexcli/pokeapi v0.0.0 => "./pokeapi"
replace github.com/dusmcd/pokedexcli/cache v0.0.0 => "./cache"
require github.com/dusmcd/pokedexcli/pokeapi v0.0.0
require github.com/dusmcd/pokedexcli/cache v0.0.0
