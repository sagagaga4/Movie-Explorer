export class YoutubeEmbed extends HTMLElement {
    // Watch for changes to data-url attribute
    static get observedAttributes() {
        return ["data-url"];
    }

    // Handle attribute changes
    attributeChangedCallback(prop, value) {
        if(prop === "data-url") {
            const url = this.dataset.url;
            console.log(url);
            const VideoId = url.substring(url.indexOf("?v") + 3);
            this.innerHTML = `
            <iframe width="560" height="315"
                src="https://www.youtube.com/embed/${VideoId}"
                title="YouTube video player" frameborder="0" allow="accelerometer;
                autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture;
                web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen>
            </iframe>
            `
        }
    }
}
// Register the custom element
customElements.define("youtube-embed", YoutubeEmbed);