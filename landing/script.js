document.addEventListener('DOMContentLoaded', function() {
    const heroSlider = document.querySelector('.slider');
    if (heroSlider) {
        const slides = heroSlider.querySelectorAll('.slide');
        const dots = heroSlider.querySelectorAll('.slider-dot');
        let currentSlide = 0;
        let slideInterval;

        function goToSlide(index) {
            if (index >= slides.length) index = 0;
            if (index < 0) index = slides.length - 1;
            
            slides[currentSlide].classList.remove('active');
            dots[currentSlide].classList.remove('active');
            
            currentSlide = index;
            
            slides[currentSlide].classList.add('active');
            dots[currentSlide].classList.add('active');
        }

        function nextSlide() {
            goToSlide(currentSlide + 1);
        }

        function startSlider() {
            slideInterval = setInterval(nextSlide, 2500);
        }

        function pauseSlider() {
            clearInterval(slideInterval);
        }

        dots.forEach((dot, index) => {
            dot.addEventListener('click', () => {
                pauseSlider();
                goToSlide(index);
                startSlider();
            });
        });

        startSlider();
    }
    
    const hamburger = document.querySelector('.hamburger');
    const nav = document.querySelector('.nav');
    const spans = hamburger.querySelectorAll('span');
    
    hamburger.addEventListener('click', function() {
        nav.classList.toggle('active');
        
        if (nav.classList.contains('active')) {
            spans[0].style.transform = 'rotate(45deg) translate(5px, 5px)';
            spans[1].style.opacity = '0';
            spans[2].style.transform = 'rotate(-45deg) translate(5px, -5px)';
        } else {
            spans[0].style.transform = 'none';
            spans[1].style.opacity = '1';
            spans[2].style.transform = 'none';
        }
    });
    
    const contactForm = document.getElementById('contact-form');
    if (contactForm) {
        contactForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const name = document.getElementById('name').value;
            const email = document.getElementById('email').value;
            const subject = document.getElementById('subject').options[document.getElementById('subject').selectedIndex].text;
            const message = document.getElementById('message').value;
            
            const mailtoLink = 'mailto:contact@chesso.org?subject=' + encodeURIComponent(subject) + 
                '&body=' + encodeURIComponent('Name: ' + name + '\r\n\r\nEmail: ' + email + '\r\n\r\nMessage:\r\n' + message);
            
            window.location.href = mailtoLink;
        });
    }
    
    window.addEventListener('scroll', function() {
        const header = document.querySelector('.header');
        if (window.scrollY > 50) {
            header.classList.add('scrolled');
        } else {
            header.classList.remove('scrolled');
        }
    });
});
