function ScoreTable(anchor) {
    if (anchor instanceof HTMLElement) { 
        this.anchor = anchor;
    }

    var selfScoreTable = this;

    this.render = function(map) {
        if (map instanceof Object) {
            selfScoreTable.map = map;

            const fragment = document.createDocumentFragment();

            const ul = document.createElement('ul');
            ul.classList.add('scores');

            for (var key in map) {
                if (map.hasOwnProperty(key)) {
                    const li = document.createElement('li');

                    const bullet = document.createElement('span');
                    bullet.innerHTML = '&#9658';
                    bullet.style.color = map[key].color;   
                    bullet.style.fontSize = '1rem';

                    const text = document.createElement('span');

                    text.id = map[key].id + '-score';

                    text.innerText = `${map[key].name}: ${map[key].score}`;
                    li.appendChild(bullet);
                    li.appendChild(text);
                    
                    ul.appendChild(li);
                }
              }

            fragment.appendChild(ul);

            this.anchor.innerHTML = '';

            this.anchor.appendChild(ul);
        }
    }
}