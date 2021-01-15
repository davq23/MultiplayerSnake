

export default class PlayerControl {
    constructor(anchor) {
        if (anchor instanceof HTMLElement) {
            this.anchor = anchor;   
            
            
            this.buttonUp = document.createElement('button');
            this.buttonUp.innerHTML = '&#8593;';
            this.buttonUp.classList.add('control-up');


            this.buttonDown = document.createElement('button');
            this.buttonDown.innerHTML = '&#8595;';
            this.buttonDown.classList.add('control-down');

            this.buttonLeft = document.createElement('button');
            this.buttonLeft.innerHTML = '&#8592;';
            this.buttonLeft.classList.add('control-left');


            this.buttonRight = document.createElement('button');
            this.buttonRight.innerHTML = '&#8594;';
            this.buttonRight.classList.add('control-right');

        }
    }

    setEvents(upEvent, downEvent, leftEvent, rightEvent) {
        this.buttonUp.onclick = upEvent;
        this.buttonDown.onclick = downEvent;
        this.buttonLeft.onclick = leftEvent;
        this.buttonRight.onclick = rightEvent;

        window.addEventListener('keydown', function (event)  {
            switch (event.key) {
                case 'w':
                case 'W':    
                    upEvent(); 
                break;       
                case 's':
                case 'S':    
                    downEvent();   
                break;     
                case 'a':
                case 'A':    
                    leftEvent();  
                break;      
                case 'd':
                case 'D':    
                    rightEvent();    
                break;
                default:
                    break;
            }
        })
    }

    render() {
        const fragment = document.createDocumentFragment();

        fragment.appendChild(this.buttonUp);
        fragment.appendChild(this.buttonDown);
        fragment.appendChild(this.buttonLeft);
        fragment.appendChild(this.buttonRight);

        const div = document.createElement('div');
        div.classList.add('control');
        div.appendChild(fragment);
        this.anchor.appendChild(div);
    }
 
}