
//Offset & Scale control
//Event listeners
function loadPageControls() {
  const viewer = document.getElementById('contentviewer');
  const iframe = document.getElementById('contentframe');
  const zoomInButton = document.getElementById('zoom-in');
  const zoomOutButton = document.getElementById('zoom-out');
  const zoomCenterButton = document.getElementById('zoom-center');
  const cont = document.querySelector('.cont');

  let initialPinchDistance = 0;
  let drag = false;
  let scale = 1;
  let offsetX = 0;
  let offsetY = 0;

  // Zoom Buttons
  zoomInButton.addEventListener('click', () => {
    scale += 0.2;
    if (scale > 3) scale = 3;
    viewer.style.setProperty('--scale', scale);
    viewer.classList.add('zoom-animation');
    viewer.addEventListener('transitionend', () => {
      viewer.classList.remove('zoom-animation');
    }, { once: true });
    centerViewer()
  });

  zoomOutButton.addEventListener('click', () => {
    scale -= 0.2;
    if (scale < 0.2) scale = 0.6;
    viewer.style.setProperty('--scale', scale);
    viewer.classList.add('zoom-animation');
    viewer.addEventListener('transitionend', () => {
      viewer.classList.remove('zoom-animation');
    }, { once: true });
    centerViewer()
  });

  function centerViewer() {
    const viewerRect = viewer.getBoundingClientRect();
    const contRect = cont.getBoundingClientRect();
    offsetX = (contRect.width - viewerRect.width) / 2;
    offsetY = (contRect.height - viewerRect.height) / 2;
    viewer.style.setProperty('--offsetX', `${offsetX}px`);
    viewer.style.setProperty('--offsetY', `${offsetY}px`);
    viewer.classList.add('zoom-animation');
    //This will make drag animation smoother while 
    // keeping the zoom animation
    viewer.addEventListener('transitionend', () => {
      viewer.classList.remove('zoom-animation');
    }, { once: true });
  }

  //Mouse Events
  viewer.addEventListener('mousedown', (e) => {
    drag = true;
    startX = e.clientX;
    startY = e.clientY;
    e.preventDefault();
  });

  document.addEventListener('mousemove', (e) => {
    if (!drag) return;
    const dx = e.clientX - startX;
    const dy = e.clientY - startY;
    offsetX += dx;
    offsetY += dy;
    viewer.style.setProperty('--offsetX', `${offsetX}px`);
    viewer.style.setProperty('--offsetY', `${offsetY}px`);
    startX = e.clientX;
    startY = e.clientY;
  });

  document.addEventListener('mouseup', () => {
    drag = false;
  });


// Touch Events
viewer.addEventListener('touchstart', (e) => {
  drag = true;
  startX = e.touches[0].clientX;
  startY = e.touches[0].clientY;
  });

  document.addEventListener('touchmove', (e) => {
  if (!drag) return;
  e.preventDefault();
  const dx = e.touches[0].clientX - startX;
  const dy = e.touches[0].clientY - startY;
  offsetX += dx;
  offsetY += dy;
  viewer.style.setProperty('--offsetX', `${offsetX}px`);
  viewer.style.setProperty('--offsetY', `${offsetY}px`);
  startX = e.touches[0].clientX;
  startY = e.touches[0].clientY;
  });

  document.addEventListener('touchend', () => {
  drag = false;
  });
/*
  viewer.addEventListener('gesturestart', (e) => {
    initialPinchDistance = e.scale;
    e.preventDefault();
  });

  viewer.addEventListener('gesturechange', (e) => {
    const newScale = scale * (e.scale / initialPinchDistance);
    if (newScale >= 0.2) {
      scale = newScale;
      viewer.style.setProperty('--scale', scale);
      centerViewer();
    }
    e.preventDefault();
  });

  viewer.addEventListener('gestureend', () => {
    initialPinchDistance = 0;
  });

*/

  // Recenter on window resize
  window.addEventListener('resize', centerViewer);

  // Recenter on window load
  //window.addEventListener('load', centerViewer);
  centerViewer();

  zoomCenterButton.addEventListener('click', () => {
    centerViewer();
  });
}

//Prep the page
//Add wrappers and buttons
function pagePrep() {
var bodyContent = document.body.innerHTML;
document.body.innerHTML = `
<div class="cont">
<div id="contentviewer">
  <div id="contentframe">
    ${bodyContent}
  </div>
</div>
</div>
`;

var zoomButtonsDiv = document.createElement('div');
zoomButtonsDiv.className = 'zoom-buttons';
zoomButtonsDiv.innerHTML = `
<div>
<a id="zoom-out"><svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512"><!--! Font Awesome Free 6.4.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. --><path d="M432 256c0 17.7-14.3 32-32 32L48 288c-17.7 0-32-14.3-32-32s14.3-32 32-32l352 0c17.7 0 32 14.3 32 32z"/></svg></a>
<a id="zoom-center"><svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512"><!--! Font Awesome Free 6.4.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. --><path d="M32 32C14.3 32 0 46.3 0 64v96c0 17.7 14.3 32 32 32s32-14.3 32-32V96h64c17.7 0 32-14.3 32-32s-14.3-32-32-32H32zM64 352c0-17.7-14.3-32-32-32s-32 14.3-32 32v96c0 17.7 14.3 32 32 32h96c17.7 0 32-14.3 32-32s-14.3-32-32-32H64V352zM320 32c-17.7 0-32 14.3-32 32s14.3 32 32 32h64v64c0 17.7 14.3 32 32 32s32-14.3 32-32V64c0-17.7-14.3-32-32-32H320zM448 352c0-17.7-14.3-32-32-32s-32 14.3-32 32v64H320c-17.7 0-32 14.3-32 32s14.3 32 32 32h96c17.7 0 32-14.3 32-32V352z"/></svg></a>
<a id="zoom-in"><svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512"><!--! Font Awesome Free 6.4.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. --><path d="M256 80c0-17.7-14.3-32-32-32s-32 14.3-32 32V224H48c-17.7 0-32 14.3-32 32s14.3 32 32 32H192V432c0 17.7 14.3 32 32 32s32-14.3 32-32V288H400c17.7 0 32-14.3 32-32s-14.3-32-32-32H256V80z"/></svg></a>
</div>
`;

document.body.appendChild(zoomButtonsDiv);
loadPageControls();
}

// Prep page on window load
window.addEventListener('load', pagePrep);
