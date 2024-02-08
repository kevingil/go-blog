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


// Home card animations
const homeCardObserver = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
        if (entry.isIntersecting) {
            entry.target.classList.add('animate-card-home');
        } else {
            entry.target.classList.remove('animate-card-home');
        }
    });
});

const observeHomeCards = (elements) => {
    elements.forEach((element) => {
        element.classList.add('animate-card-home');
        homeCardObserver.observe(element);
    });
};


// Home text animations
const textHomeObserver = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
        if (entry.isIntersecting) {
            entry.target.classList.add('animate-down');
        } else {
            entry.target.classList.remove('animate-down');
        }
    });
});

const observeHomeText = (elements) => {
    elements.forEach((element) => {
        textHomeObserver.observe(element);
    });
};

// Combined onLoad function
htmx.onLoad(function (elt) {
    // Text animations
    const hiddenElements = elt.querySelectorAll('.hide');
    observeElements(hiddenElements);

    // Card animations
    const hiddenCards = elt.querySelectorAll('.hide-card');
    observeCards(hiddenCards);

    // Home animations
    const hiddenHomeCards = elt.querySelectorAll('.hide-card-home');
    observeHomeCards(hiddenHomeCards);

    // Text home animations
    const hiddenHomeText = elt.querySelectorAll('.hide-down');
    observeHomeText(hiddenHomeText);

});




// HTMX on hx-swap, scroll to top, no animation
/*
document.addEventListener('htmx:afterSwap', function (event) {
    window.scrollTo(0, 0);
});
*/
