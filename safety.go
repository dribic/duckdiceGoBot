/*
 Copyright (C) 2024 Dejan Ribič

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import "math"

func LabSafety(startBet float64, inTable []float64) []float64 {
	var newTable []float64

	for idx, element := range inTable {
		if idx <= len(inTable)-3 {
			for element >= startBet {
				newTable = append(newTable, startBet)
				element -= startBet
				element = math.Round(element*1000) / 1000
			}
		} else {
			for element > 0 {
				if element >= 2*startBet {
					newTable = append(newTable, 2*startBet)
					element -= 2 * startBet
					element = math.Round(element*1000) / 1000
				} else {
					newTable = append(newTable, startBet)
					element = 0
				}
			}
		}
	}

	return newTable
}
