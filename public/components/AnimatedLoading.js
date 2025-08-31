class AnimatedLoading extends HTMLElement{
    constructor(){
        super();
    }

    connectedCallback(){
        const elements = this.dataset.elements;
        const width = this.dataset.width;
        const height = this.dataset.height;

        for(let i = 0; i<elements; i++){
            const wrapper = document.createElement("div");
            wrapper.style.width = width;
            wrapper.style.height = height;
            wrapper.style.margin = "10px";
            wrapper.style.color= "white";
            wrapper.style.display = "inline-block"
            this.appendChild(wrapper);
        }
    }
}

customElements.define("animated-loading", AnimatedLoading);

