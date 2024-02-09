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

// Card stack animations
const cardStackObserver = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
        if (entry.isIntersecting) {
            entry.target.classList.add('animate-card-stack');
        } else {
            entry.target.classList.remove('animate-card-stack');
        }
    });
});

const observeCardStack = (elements) => {
    elements.forEach((element) => {
        cardStackObserver.observe(element);
    });
};


// Card animations
const c1Observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
        if (entry.isIntersecting) {
            entry.target.classList.add('animate-c1');
        } else {
            entry.target.classList.remove('animate-c1');
        }
    });
});

const observeC1 = (elements) => {
    elements.forEach((element) => {
        c1Observer.observe(element);
    });
};

// Card animations
const c2Observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
        if (entry.isIntersecting) {
            entry.target.classList.add('animate-c2');
        } else {
            entry.target.classList.remove('animate-c2');
        }
    });
});

const observeC2 = (elements) => {
    elements.forEach((element) => {
        c2Observer.observe(element);
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
htmx.onLoad(function (elt) {
    // Text animations
    const hiddenElements = elt.querySelectorAll('.hide');
    observeElements(hiddenElements);

    // Card animations
    const hiddenCards = elt.querySelectorAll('.hide-card');
    observeCards(hiddenCards);


    // Card animations
    const hiddenCardStack = elt.querySelectorAll('.hide-cards');
    observeCardStack(hiddenCardStack);

    // C1 animations
    const hiddenC1 = elt.querySelectorAll('.fade-c1');
    observeC1(hiddenC1);

    // C2 animations
    const hiddenC2 = elt.querySelectorAll('.fade-c2');
    observeC2(hiddenC2);

    // Home animations
    const hiddenHomeCards = elt.querySelectorAll('.hide-card-home');
    hiddenHomeCards.forEach(element => {
        element.classList.add('animate-card-home');
    });

    // Text home animations
    const hiddenHomeText = elt.querySelectorAll('.hide-down');
    hiddenHomeText.forEach(element => {
        element.classList.add('animate-down');
    });
});




// HTMX on hx-swap, scroll to top, no animation
/*
document.addEventListener('htmx:afterSwap', function (event) {
    window.scrollTo(0, 0);
});
*/
