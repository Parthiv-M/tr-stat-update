package types

import "time"

type AccessToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type Credits struct {
	Type             string `json:"type"`
	Total            int    `json:"total"`
	Index            int    `json:"index"`
	NValues          int    `json:"nValues"`
	TrpCreditPackage []struct {
		PackageID    int       `json:"packageId"`
		ProductID    int       `json:"productId"`
		PurchaseDate time.Time `json:"purchaseDate"`
		Active       bool      `json:"active"`
		UserID       int       `json:"userId"`
		UserName     string    `json:"userName"`
		FreeOfCharge bool      `json:"freeOfCharge"`
		OrigUserID   int       `json:"origUserId"`
		Product      struct {
			ProductID   int    `json:"productId"`
			CreditType  string `json:"creditType"`
			NrOfCredits int    `json:"nrOfCredits"`
			Label       string `json:"label"`
			Shareable   bool   `json:"shareable"`
			Status      int    `json:"status"`
		} `json:"product,omitempty"`
		Balance float64 `json:"balance"`
	} `json:"trpCreditPackage"`
	OverallBalance float64 `json:"overallBalance"`
}

type Page struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Normalized []struct {
			Fromencoded bool   `json:"fromencoded"`
			From        string `json:"from"`
			To          string `json:"to"`
		} `json:"normalized"`
		Pages []struct {
			Pageid    int    `json:"pageid"`
			Ns        int    `json:"ns"`
			Title     string `json:"title"`
			Revisions []struct {
				Slots struct {
					Main struct {
						Contentmodel  string `json:"contentmodel"`
						Contentformat string `json:"contentformat"`
						Content       string `json:"content"`
					} `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

type LoginTokens struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Tokens struct {
			Logintoken string `json:"logintoken"`
		} `json:"tokens"`
	} `json:"query"`
}

type CSRFTokens struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Tokens struct {
			CSRFtoken string `json:"csrftoken"`
		} `json:"tokens"`
	} `json:"query"`
}

type LoginData struct {
	Login struct {
		Result     string `json:"result"`
		Lguserid   int    `json:"lguserid"`
		Lgusername string `json:"lgusername"`
	} `json:"login"`
}

type WikiPageData struct {
	License     string `json:"license"`
	Description struct {
		En string `json:"en"`
	} `json:"description"`
	Sources string `json:"sources"`
	Schema  struct {
		Fields []struct {
			Name  string `json:"name"`
			Type  string `json:"type"`
			Title struct {
				En string `json:"en"`
			} `json:"title"`
		} `json:"fields"`
	} `json:"schema"`
	Data [][]interface{} `json:"data"`
}

type EditResponse struct {
	Edit struct {
		Result       string    `json:"result"`
		Pageid       int       `json:"pageid"`
		Title        string    `json:"title"`
		Contentmodel string    `json:"contentmodel"`
		Oldrevid     int       `json:"oldrevid"`
		Newrevid     int       `json:"newrevid"`
		Newtimestamp time.Time `json:"newtimestamp"`
		Watched      bool      `json:"watched"`
	} `json:"edit"`
}
