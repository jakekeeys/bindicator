package collection

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type BinCollectionDates struct {
	Household time.Time
	Recycling time.Time
	Food time.Time
	Garden time.Time
}

func GetNext(ctx context.Context, debug bool, url, postcode, number string) (*BinCollectionDates, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	if debug {
		ctx, cancel = chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
		defer cancel()
	}

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	postcodeSearchSelector := `//input[@name="CollectionDayLookup2$TextBox_PostCode"]`
	postcodeSearchSubmit := `//input[@name="CollectionDayLookup2$Button_PostCodeSearch"]`
	addressSelector := `//select[@name="CollectionDayLookup2$DropDownList_Addresses"]`
	addressSelectorSubmit := `//input[@name="CollectionDayLookup2$Button_SelectAddress"]`

	var household, recycling, food, garden string
	err := chromedp.Run(
		ctx,

		// Load the bin collections page and search by postcode
		chromedp.Navigate(url),
		chromedp.WaitVisible(postcodeSearchSelector),
		chromedp.SendKeys(postcodeSearchSelector, postcode),
		chromedp.Click(postcodeSearchSubmit),

		// Select the correct address
		chromedp.WaitNotPresent(addressSelector),
		chromedp.SendKeys(addressSelector, number),
		chromedp.Click(addressSelectorSubmit),

		// Load the collection data
		chromedp.WaitNotPresent(`//h4[contains(., 'Collection Dates')]`),
		chromedp.Text("#CollectionDayLookup2_Label_HouseholdWaste_Date", &household, chromedp.ByQuery),
		chromedp.Text("#CollectionDayLookup2_Label_RecyclingWaste_Date", &recycling, chromedp.ByQuery),
		chromedp.Text("#CollectionDayLookup2_Label_FoodWaste_Date", &food, chromedp.ByQuery),
		chromedp.Text("#CollectionDayLookup2_Label_GardenWaste_Date", &garden, chromedp.ByQuery),

	)
	if err != nil {
		return nil, fmt.Errorf("could not run chromedp: %w", err)
	}

	householdTS, err := getTimeFromCollectionDateString(household)
	if err != nil {
		return nil, fmt.Errorf("error parsing collection date: %w", err)
	}

	recyclingTS, err := getTimeFromCollectionDateString(recycling)
	if err != nil {
		return nil, fmt.Errorf("error parsing collection date: %w", err)
	}

	foodTS, err := getTimeFromCollectionDateString(food)
	if err != nil {
		return nil, fmt.Errorf("error parsing collection date: %w", err)
	}

	gardenTS, err := getTimeFromCollectionDateString(garden)
	if err != nil {
		return nil, fmt.Errorf("error parsing collection date: %w", err)
	}

	return &BinCollectionDates{
		Household: householdTS.Local(),
		Recycling: recyclingTS.Local(),
		Food:      foodTS.Local(),
		Garden:    gardenTS.Local(),
	}, nil
}

func getTimeFromCollectionDateString(s string) (*time.Time, error) {
	var ts time.Time
	if strings.ToLower(s) == "today" {
		ts = time.Now().Truncate(time.Hour * 24)
	} else {
		var err error
		ts, err = time.Parse("Monday 02/01/06", s)
		if err != nil {
			return nil, fmt.Errorf("could not parse collection timestamp: %w", err)
		}
	}

	return &ts, nil
}