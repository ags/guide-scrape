#!/bin/bash

set -euo pipefail

function main {
  local regions=(auckland brisbane melbourne sydney wellington)
  local types=(restaurant bar cafe pub)

  for region in ${regions[@]}; do
    for t in ${types[@]}; do
      concreteplayground $region $t > "concreteplayground-${region}-${t}.csv"
    done
  done
}

main
