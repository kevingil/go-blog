// Header styling
window.addEventListener('scroll', function () {
    var scrollPosition = window.scrollY;
    var element = document.querySelector('.scrollfade');
    if (scrollPosition > 30) {
        element.classList.add('scrolled');
    } else {
        element.classList.remove('scrolled');
    }
});

// Text animations (yes, I watched the fireship video)
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

// Normal hidden elements
const hiddenElements = document.querySelectorAll('.hide');
observeElements(hiddenElements);

// Animate newly loaded elements
htmx.onLoad(function(elt) {
    const newHiddenElements = elt.querySelectorAll('.hide');
    observeElements(newHiddenElements);
});


// Card animations

// Animation observer
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

// Normal hidden elements
const hiddenCards = document.querySelectorAll('.hide-card');
observeCards(hiddenCards);

// Animate newly loaded elements
htmx.onLoad(function(elt) {
    const newHiddenCards = elt.querySelectorAll('.hide-card');
    observeCards(newHiddenCards);
});



// HTMX on hx-swap, scroll to top, no animation
document.addEventListener('htmx:afterSwap', function (event) {
    window.scrollTo(0, 0);
});
