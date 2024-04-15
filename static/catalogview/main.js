

function init() {
  // Must include width and height of content in px
  const catalog = new Catalog('921', '1196')
}
document.addEventListener("DOMContentLoaded", () => {
  init()
});

class Catalog {
  constructor(width, height) {
    this.debug = false
    this.content_width = width
    this.content_height = height
    this.wrap()
    this.startX = null
    this.startY = null
    this.initialPinchDistance = 0
    this.drag = false
    this.dx = 0
    this.dy = 0
    this.offsetX = 0
    this.offsetY = 0
    this.content = document.getElementById('content')
    this._content = this.content.getBoundingClientRect()
    this.viewer = document.getElementById('viewer')
    this._viewer = this.viewer.getBoundingClientRect()
    this._viewer = this.viewer.getBoundingClientRect()
    this.listen()
    this.update_scale((this._viewer.width / this.content_width))
    this.center_offset()

  }

  center_offset() {

    this.offsetX = ((this._viewer.width) - this.content_width) / 2;
    this.offsetY = ((this._viewer.height) - this.content_height) / 2;
    this.content.style.setProperty('--offsetX', `${this.offsetX}px`);
    this.content.style.setProperty('--offsetY', `${this.offsetY}px`);
    this.log(`center offset`);
  }


  update_scale(s) {
    this.scale = s
    this._content = this.content.getBoundingClientRect();
    this._viewer = this.viewer.getBoundingClientRect();

    this.content.style.setProperty('--scale', this.scale);
    this.offsetX = ((this._viewer.width) - this._content.width) / 2;
    this.offsetY = ((this._viewer.height) - this._content.height) / 2;

    // Update CSS properties
    this.content.style.setProperty('--offsetX', `${this.offsetX}px`);
    this.content.style.setProperty('--offsetY', `${this.offsetY}px`);

    //Animate onclick only
    this.content.classList.add('zoom-animation')
    this.content.addEventListener('transitionend', () => {
      this.content.classList.remove('zoom-animation')
    }, { once: true })

    this.log('update scale');
  }

  //Activate event listeners
  listen() {

    document.getElementById('zoom-in').addEventListener('click', () => {
      let scale = this.scale + 0.1
      if (scale > 1.5) scale = 1.51

      this.update_scale(scale)
      this.center_offset()

    })

    document.getElementById('zoom-out').addEventListener('click', () => {
      let scale = this.scale - 0.1
      if (scale < 0.4) scale = 0.39

      this.update_scale(scale)
      this.center_offset()

    })

    //Center button
    document.getElementById('zoom-center').addEventListener('click', () => {
      let scale = (this.viewer.getBoundingClientRect().width / this.content_width)
      if (scale > 1) scale = 1
      console.log(`DEBUG ${this.viewer.getBoundingClientRect().width} + ${this.content.getBoundingClientRect().width} + ${scale}`)
      this.update_scale(scale)
      this.center_offset()
    })

    // Recenter on window resize
    window.addEventListener('resize', () => {

      let scale = (this.viewer.getBoundingClientRect().width / this.content_width)
      if (scale > 1) scale = 1
      console.log(`DEBUG ${this.viewer.getBoundingClientRect().width} + ${this.content_width} + ${scale}`)

      this.update_scale(scale)
      this.center_offset()
    })


    //MOUSE MOVE
    this.content.addEventListener('mousedown', (e) => {
      this.drag = true;
      this.startX = e.clientX;
      this.startY = e.clientY;
      e.preventDefault()
    })

    document.addEventListener('mousemove', (e) => {
      if (!this.drag) return;
      this.dx = e.clientX - this.startX;
      this.dy = e.clientY - this.startY;
      this.offsetX += this.dx;
      this.offsetY += this.dy;
      this.content.style.setProperty('--offsetX', `${this.offsetX}px`)
      this.content.style.setProperty('--offsetY', `${this.offsetY}px`)
      this.startX = e.clientX;
      this.startY = e.clientY;
    })

    document.addEventListener('mouseup', () => {
      this.drag = false;
      this.log('mouseup')
    })


    // TOUCH MOVE
    this.content.addEventListener('touchstart', (e) => {
      this.drag = true;
      this.startX = e.touches[0].clientX;
      this.startY = e.touches[0].clientY;
      if (e.target.tagName.toLowerCase() !== 'a') {
        e.preventDefault()

        console.log("Content Position X, Y: " + this._content.top + ", " + this._content.left)
      }
    })

    document.addEventListener('touchmove', (e) => {
      if (!this.drag) return;
      this.dx = e.touches[0].clientX - this.startX;
      this.dy = e.touches[0].clientY - this.startY;
      this.offsetX += this.dx;
      this.offsetY += this.dy;
      this.content.style.setProperty('--offsetX', `${this.offsetX}px`)
      this.content.style.setProperty('--offsetY', `${this.offsetY}px`)
      this.startX = e.touches[0].clientX;
      this.startY = e.touches[0].clientY;
      if (e.target.tagName.toLowerCase() !== 'a') {
        e.preventDefault()
        e.stopPropagation()
      }
    })

    document.addEventListener('touchend', () => {
      this.drag = false;
    })
  }



  log(e = '') {
    if (this.debug) {
      console.log(`========================================`)
      console.log(`EVENT: ${e}`)
      console.log(`V RECT: ${this._viewer.width}`)
      console.log(`C RECT: ${this._content.width}, ${this._content.height}`)
      console.log(`OFFSET x: ${this.offsetX} x: ${this.offsetY}`)
      console.log(`Scale: ${this.scale}`)
      console.log(`========================================`)
    }
  }


  //Wrap the page in a viewer and content div
  wrap() {
    var bodyContent = document.body.innerHTML;
    document.body.innerHTML = `
  <div id="viewer">
  <div id="content" style="width:${this.content_width}px; height: ${this.content_height}px; ">
      ${bodyContent}
  </div>
  </div>
  `;

    var zoomButtonsDiv = document.createElement('div')
    zoomButtonsDiv.className = 'zoom-buttons';
    zoomButtonsDiv.innerHTML = `
  <div>
    <a id="zoom-out"><svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512">
      <path d="M432 256c0 17.7-14.3 32-32 32L48 288c-17.7 0-32-14.3-32-32s14.3-32 32-32l352 0c17.7 0 32 14.3 32 32z"/></svg>
    </a>
    <a id="zoom-center"><svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512">
      <path d="M32 32C14.3 32 0 46.3 0 64v96c0 17.7 14.3 32 32 32s32-14.3 32-32V96h64c17.7 0 32-14.3 32-32s-14.3-32-32-32H32zM64 352c0-17.7-14.3-32-32-32s-32 14.3-32 32v96c0 17.7 14.3 32 32 32h96c17.7 0 32-14.3 32-32s-14.3-32-32-32H64V352zM320 32c-17.7 0-32 14.3-32 32s14.3 32 32 32h64v64c0 17.7 14.3 32 32 32s32-14.3 32-32V64c0-17.7-14.3-32-32-32H320zM448 352c0-17.7-14.3-32-32-32s-32 14.3-32 32v64H320c-17.7 0-32 14.3-32 32s14.3 32 32 32h96c17.7 0 32-14.3 32-32V352z"/></svg>
    </a>
    <a id="zoom-in">
      <svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512"><path d="M256 80c0-17.7-14.3-32-32-32s-32 14.3-32 32V224H48c-17.7 0-32 14.3-32 32s14.3 32 32 32H192V432c0 17.7 14.3 32 32 32s32-14.3 32-32V288H400c17.7 0 32-14.3 32-32s-14.3-32-32-32H256V80z"/></svg>
    </a>
  </div>
  `;

    document.body.appendChild(zoomButtonsDiv)
  }

}
