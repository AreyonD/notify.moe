'use strict'

let Promise = require('bluebird')
let chalk = require('chalk')
let fs = Promise.promisifyAll(require('fs'))

let updateGenres = () => {
	console.log(chalk.yellow('✖'), 'Updating genre cache...')

	fs.readFileAsync('pages/anime/genres/genres.txt', 'utf8').then(genreText => {
		let genreList = genreText.split('\n')

		genreList.forEach(genre => {
			genre = arn.fixGenre(genre)
			let genreSearch = `;${genre};`

			arn.cacheLimiter.removeTokens(1, () => {
				let animeList = []

				return arn.forEach('Anime', anime => {
					if(!anime.watching)
						return

					if(!anime.genres)
						return

					if((';' + anime.genres.map(arn.fixGenre).join(';') + ';').indexOf(genreSearch) === -1)
						return

					animeList.push(anime)
				}).then(() => {
					animeList.sort((a, b) => {
						if(a.watching === b.watching) {
							if(!a.startDate)
								return 1

							if(!b.startDate)
								return -1

							return a.startDate > b.startDate ? -1 : 1
						} else if(!a.watching) {
							return 1
						} else if(!b.watching) {
							return -1
						} else {
							return a.watching > b.watching ? -1 : 1
						}
					})

					console.log(chalk.green('✔'), `Updated genre ${chalk.yellow(genre)} (${animeList.length} anime)`)
					
					return arn.set('Genres', genre, {
						genre,
						animeList
					})
				})
			})
		})
	})
}

arn.repeatedly(30 * 60, updateGenres)