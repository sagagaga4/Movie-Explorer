import { API } from "../services/API.js";
import { CollectionPage } from "./CollectionPage.js";

export default class WatchlistPage extends CollectionPage {

    constructor() {
        super(API.getWatchlist, "Movie Watchlist")
    }

}
customElements.define("watchlist-page", WatchlistPage);
