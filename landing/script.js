document.addEventListener('DOMContentLoaded', function() {
    // Slider functionality
    const slider = document.querySelector('.slider');
    if (slider) {
        const slides = slider.querySelectorAll('.slide');
        const dots = slider.querySelectorAll('.slider-dot');
        const prevArrow = slider.querySelector('.slider-arrow.prev');
        const nextArrow = slider.querySelector('.slider-arrow.next');
        let currentSlide = 0;
        let slideInterval;
        let isTransitioning = false;

        function showSlide(index) {
            if (isTransitioning) return;
            isTransitioning = true;
            
            // Ensure index is within bounds
            if (index >= slides.length) index = 0;
            if (index < 0) index = slides.length - 1;
            
            // Remove active classes
            slides.forEach(slide => slide.classList.remove('active'));
            dots.forEach(dot => dot.classList.remove('active'));
            
            // Add active classes
            slides[index].classList.add('active');
            dots[index].classList.add('active');
            
            currentSlide = index;
            
            // Allow transitions after animation completes (matches CSS transition)
            setTimeout(() => {
                isTransitioning = false;
            }, 800);
        }

        function nextSlide() {
            if (isTransitioning) return;
            showSlide(currentSlide + 1);
        }

        function prevSlide() {
            if (isTransitioning) return;
            showSlide(currentSlide - 1);
        }

        function startAutoplay() {
            if (slideInterval) {
                clearInterval(slideInterval);
            }
            slideInterval = setInterval(() => {
                if (!isTransitioning) {
                    nextSlide();
                }
            }, 6000);
        }

        function stopAutoplay() {
            if (slideInterval) {
                clearInterval(slideInterval);
                slideInterval = null;
            }
        }

        // Arrow navigation
        if (nextArrow) {
            nextArrow.addEventListener('click', () => {
                stopAutoplay();
                nextSlide();
                setTimeout(startAutoplay, 1000);
            });
        }

        if (prevArrow) {
            prevArrow.addEventListener('click', () => {
                stopAutoplay();
                prevSlide();
                setTimeout(startAutoplay, 1000);
            });
        }

        // Dot navigation
        dots.forEach((dot, index) => {
            dot.addEventListener('click', () => {
                if (index !== currentSlide) {
                    stopAutoplay();
                    showSlide(index);
                    setTimeout(startAutoplay, 1000);
                }
            });
        });

        // Pause on hover
        slider.addEventListener('mouseenter', stopAutoplay);
        slider.addEventListener('mouseleave', () => {
            if (!slideInterval) {
                startAutoplay();
            }
        });

        // Initialize
        showSlide(0);
        startAutoplay();
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
