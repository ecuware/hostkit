// Extra JavaScript for HostKit documentation

// Add copy button to code blocks
document.addEventListener('DOMContentLoaded', function() {
  // Smooth scroll for anchor links
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
      e.preventDefault();
      const target = document.querySelector(this.getAttribute('href'));
      if (target) {
        target.scrollIntoView({
          behavior: 'smooth',
          block: 'start'
        });
      }
    });
  });

  // Add external link indicators
  document.querySelectorAll('a[href^="http"]').forEach(link => {
    if (!link.href.includes(window.location.hostname)) {
      link.setAttribute('target', '_blank');
      link.setAttribute('rel', 'noopener noreferrer');
      
      // Add external link icon if not already present
      if (!link.querySelector('.external-link-icon')) {
        const icon = document.createElement('span');
        icon.className = 'external-link-icon';
        icon.innerHTML = ' ↗';
        icon.style.fontSize = '0.8em';
        link.appendChild(icon);
      }
    }
  });
});

// Console welcome message
console.log('%cHostKit Documentation', 'font-size: 24px; font-weight: bold; color: #4051b5;');
console.log('%cAll-in-one hosting server management toolkit', 'font-size: 14px; color: #666;');
console.log('%chttps://github.com/ecuware/hostkit', 'font-size: 12px; color: #526cfe;');
