package companies

import (
	"github.com/aerogo/aero"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"
	"github.com/animenotifier/notify.moe/utils/infinitescroll"
)

const maxPopularCompanies = 50

// Popular renders the best companies.
func Popular(ctx *aero.Context) string {
	user := utils.GetUser(ctx)
	index, _ := ctx.GetInt("index")

	// Fetch all eligible companies
	allCompanies := fetchAll()

	// Sort the companies by popularity
	arn.SortCompaniesPopularFirst(allCompanies)

	// Slice the part that we need
	companies := allCompanies[index:]

	if len(companies) > maxPopularCompanies {
		companies = companies[:maxPopularCompanies]
	}

	// Next index
	nextIndex := infinitescroll.NextIndex(ctx, len(allCompanies), maxPopularCompanies, index)

	// In case we're scrolling, send companies only (without the page frame)
	if index > 0 {
		return ctx.HTML(components.PopularCompaniesScrollable(companies, user))
	}

	// Otherwise, send the full page
	return ctx.HTML(components.PopularCompanies(companies, nextIndex, user))
}
