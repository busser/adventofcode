package busser

import (
	"testing"
)

func TestTileRotate(t *testing.T) {
	/*
		original:
			###....###
			##......##
			#...##...#
			#.......##
			##..#..###
			##....####
			#.....####
			.......###
			####....##
			..###....#

		rotated:
			.#.#######
			.#..##..##
			##.......#
			##........
			#....#.#..
			.......#..
			...##.....
			..####...#
			.######.##
			##########
	*/
	f, e := pixelFull, pixelEmpty
	original := tile{
		content: [tileSize][tileSize]pixel{
			{f, f, f, e, e, e, e, f, f, f},
			{f, f, e, e, e, e, e, e, f, f},
			{f, e, e, e, f, f, e, e, e, f},
			{f, e, e, e, e, e, e, e, f, f},
			{f, f, e, e, f, e, e, f, f, f},
			{f, f, e, e, e, e, f, f, f, f},
			{f, e, e, e, e, e, f, f, f, f},
			{e, e, e, e, e, e, e, f, f, f},
			{f, f, f, f, e, e, e, e, f, f},
			{e, e, f, f, f, e, e, e, e, f},
		},
	}
	rotated := tile{
		content: [tileSize][tileSize]pixel{
			{e, f, e, f, f, f, f, f, f, f},
			{e, f, e, e, f, f, e, e, f, f},
			{f, f, e, e, e, e, e, e, e, f},
			{f, f, e, e, e, e, e, e, e, e},
			{f, e, e, e, e, f, e, f, e, e},
			{e, e, e, e, e, e, e, f, e, e},
			{e, e, e, f, f, e, e, e, e, e},
			{e, e, f, f, f, f, e, e, e, f},
			{e, f, f, f, f, f, f, e, f, f},
			{f, f, f, f, f, f, f, f, f, f},
		},
	}

	original.rotate()

	if original != rotated {
		t.Error("rotation algorithm does not work")
	}
}

func TestTileFlip(t *testing.T) {
	/*
		original:
			###....###
			##......##
			#...##...#
			#.......##
			##..#..###
			##....####
			#.....####
			.......###
			####....##
			..###....#

		flipped:
			###....###
			##......##
			#...##...#
			##.......#
			###..#..##
			####....##
			####.....#
			###.......
			##....####
			#....###..
	*/
	f, e := pixelFull, pixelEmpty
	original := tile{
		content: [tileSize][tileSize]pixel{
			{f, f, f, e, e, e, e, f, f, f},
			{f, f, e, e, e, e, e, e, f, f},
			{f, e, e, e, f, f, e, e, e, f},
			{f, e, e, e, e, e, e, e, f, f},
			{f, f, e, e, f, e, e, f, f, f},
			{f, f, e, e, e, e, f, f, f, f},
			{f, e, e, e, e, e, f, f, f, f},
			{e, e, e, e, e, e, e, f, f, f},
			{f, f, f, f, e, e, e, e, f, f},
			{e, e, f, f, f, e, e, e, e, f},
		},
	}
	flipped := tile{
		content: [tileSize][tileSize]pixel{
			{f, f, f, e, e, e, e, f, f, f},
			{f, f, e, e, e, e, e, e, f, f},
			{f, e, e, e, f, f, e, e, e, f},
			{f, f, e, e, e, e, e, e, e, f},
			{f, f, f, e, e, f, e, e, f, f},
			{f, f, f, f, e, e, e, e, f, f},
			{f, f, f, f, e, e, e, e, e, f},
			{f, f, f, e, e, e, e, e, e, e},
			{f, f, e, e, e, e, f, f, f, f},
			{f, e, e, e, e, f, f, f, e, e},
		},
	}

	original.flip()

	if original != flipped {
		t.Error("flipping algorithm does not work")
	}
}
