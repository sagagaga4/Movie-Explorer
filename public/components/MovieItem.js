export class MovieItemComponent extends HTMLElement {
    constructor(movie) {
        super();
        this.movie = movie;
    }

    connectedCallback() {
        const url = "/movies/" + this.movie.id;
        this.innerHTML = `
                <a class ="navlink" href="#" onclick="event.preventDefault();app.Router.go('${url}')">
                    <article>
                        <img src="${this.movie.poster_url}" alt="${this.movie.title} Poster">
                        <p>${this.movie.title} (${this.movie.release_year})</p>
                    </article>
                </a>
            `
    }
}

customElements.define("movie-item", MovieItemComponent);