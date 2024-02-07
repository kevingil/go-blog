// Text animations 
const observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
        if (entry.isIntersecting) {
            entry.target.classList.add('animate');
        } else {
            entry.target.classList.remove('animate');
        }
    });
});

const observeElements = (elements) => {
    elements.forEach((element) => {
        observer.observe(element);
    });
};

// Card animations
const cardObserver = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
        if (entry.isIntersecting) {
            entry.target.classList.add('animate-card');
        } else {
            entry.target.classList.remove('animate-card');
        }
    });
});

const observeCards = (elements) => {
    elements.forEach((element) => {
        cardObserver.observe(element);
    });
};

// Combined onLoad function
htmx.onLoad(function(elt) {
    // Text animations
    const hiddenElements = elt.querySelectorAll('.hide');
    observeElements(hiddenElements);

    // Card animations
    const hiddenCards = elt.querySelectorAll('.hide-card');
    observeCards(hiddenCards);

    // Animate newly loaded elements
    const newHiddenElements = elt.querySelectorAll('.hide');
    const newHiddenCards = elt.querySelectorAll('.hide-card');

    observeElements(newHiddenElements);
    observeCards(newHiddenCards);

    // Onload animations
    const hideCardHome = elt.querySelectorAll('.hide-card-home');
    hideCardHome.forEach(card => {
        card.classList.add('animate-card-home');
    });

    const hideHome = elt.querySelectorAll('.hide-down');
    hideHome.forEach(card => {
        card.classList.add('animate-down');
    });
});




// HTMX on hx-swap, scroll to top, no animation
/*
document.addEventListener('htmx:afterSwap', function (event) {
    window.scrollTo(0, 0);
});
*/
