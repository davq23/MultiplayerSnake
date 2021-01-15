

export default class Scores {
    constructor(anchor) {
        if (anchor instanceof HTMLElement) { 
            this.anchor = anchor;
        }
    }

    render(map) {
        if (map instanceof Map) { 
            this.map = map;

            const fragment = document.createDocumentFragment();

            const ul = document.createElement('ul');
            ul.classList.add('scores');

            map.forEach(function(value, key, map) {
                const li = document.createElement('li');

                const bullet = document.createElement('span');
                bullet.innerHTML = '&#9658';
                bullet.style.color = value.color;   
                bullet.style.fontSize = '1rem';

                const text = document.createElement('span');

                text.id = value.id + '-score';

                text.innerText = `${value.name}: ${value.score}`;
                li.appendChild(bullet);
                li.appendChild(text);
                
                ul.appendChild(li);
            });

            fragment.appendChild(ul);

            this.anchor.innerHTML = '';

            this.anchor.appendChild(ul);
        }
    } 
}