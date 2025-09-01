import React, { useEffect, useState } from 'react';
import { Copy, ArrowRight, Github, AlertTriangle, Code, Zap, FileText, MessageCircle } from "lucide-react";

export function StarknetLanding() {
  const [windowWidth, setWindowWidth] = useState(1200);
  const [bannerIsAnimating, setBannerIsAnimating] = useState(true);

  useEffect(() => {
    const updateWidth = () => setWindowWidth(window.innerWidth);
    updateWidth();
    window.addEventListener('resize', updateWidth);
    return () => window.removeEventListener('resize', updateWidth);
  }, []);

  useEffect(() => {
    const timer = setTimeout(() => setBannerIsAnimating(false), 3000);
    return () => clearTimeout(timer);
  }, []);

  // Handle scrollbar visibility on scroll
  useEffect(() => {
    let scrollTimer: NodeJS.Timeout;
    
    const handleScroll = () => {
      document.body.classList.add('scrolling');
      
      // Clear existing timer
      clearTimeout(scrollTimer);
      
      // Hide scrollbar after scroll stops (after 1 second)
      scrollTimer = setTimeout(() => {
        document.body.classList.remove('scrolling');
      }, 1000);
    };

    // Add listeners to both window and the main container
    window.addEventListener('scroll', handleScroll, true);
    const mainContainer = document.querySelector('[style*="position: fixed"][style*="overflow: auto"]');
    if (mainContainer) {
      mainContainer.addEventListener('scroll', handleScroll);
    }
    
    return () => {
      window.removeEventListener('scroll', handleScroll, true);
      if (mainContainer) {
        mainContainer.removeEventListener('scroll', handleScroll);
      }
      clearTimeout(scrollTimer);
    };
  }, []);

  // Hide vocs.dev header and make full screen
  useEffect(() => {
    // Hide vocs-specific elements
    const vocsElements = [
      '.vocs_Header',
      '.vocs_Nav',
      '.vocs_TopNav'
    ];
    
    vocsElements.forEach(selector => {
      const element = document.querySelector(selector);
      if (element) {
        (element as HTMLElement).style.display = 'none';
      }
    });

    // Modify parent containers
    const containerSelectors = [
      '[data-layout="landing"]',
      'main',
      'article',
      '.vocs_Content',
      '.vocs_ContentWrapper'
    ];
    
    containerSelectors.forEach(selector => {
      const container = document.querySelector(selector);
      if (container) {
        const el = container as HTMLElement;
        el.style.padding = '0';
        el.style.margin = '0';
        el.style.maxWidth = 'none';
        el.style.width = '100%';
      }
    });

    // Cleanup function
    return () => {
      vocsElements.forEach(selector => {
        const element = document.querySelector(selector);
        if (element) {
          (element as HTMLElement).style.display = '';
        }
      });
    };
  }, []);

  const copyToClipboard = () => {
    navigator.clipboard.writeText("go get github.com/NethermindEth/starknet.go");
    
    // Create notification container
    const notification = document.createElement('div');
    notification.style.cssText = `
      position: fixed;
      top: 24px;
      right: 24px;
      background: rgba(17, 24, 39, 0.95);
      color: white;
      padding: 16px 20px;
      border-radius: 8px;
      z-index: 10001;
      font-size: 14px;
      font-weight: 500;
      font-family: ui-sans-serif, system-ui, sans-serif;
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.4), 0 10px 10px -5px rgba(0, 0, 0, 0.2);
      border: 1px solid rgba(75, 85, 99, 0.3);
      backdrop-filter: blur(12px);
      min-width: 200px;
      transform: translateX(100%);
      transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
      overflow: hidden;
    `;
    
    // Create progress bar
    const progressBar = document.createElement('div');
    progressBar.style.cssText = `
      position: absolute;
      top: 0;
      left: 0;
      height: 2px;
      background: hsl(14 97% 49%);
      width: 100%;
      transform-origin: left;
      animation: progressShrink 3000ms linear forwards;
    `;
    
    // Add keyframes for progress bar animation
    if (!document.querySelector('#progress-keyframes')) {
      const style = document.createElement('style');
      style.id = 'progress-keyframes';
      style.textContent = `
        @keyframes progressShrink {
          from { transform: scaleX(1); }
          to { transform: scaleX(0); }
        }
      `;
      document.head.appendChild(style);
    }
    
    notification.innerHTML = 'Copied to clipboard';
    notification.appendChild(progressBar);
    document.body.appendChild(notification);
    
    // Trigger slide-in animation
    setTimeout(() => {
      notification.style.transform = 'translateX(0)';
    }, 10);
    
    // Auto-hide with slide-out animation
    setTimeout(() => {
      notification.style.transform = 'translateX(100%)';
      notification.style.opacity = '0';
    }, 2800);
    
    // Remove from DOM
    setTimeout(() => {
      if (document.body.contains(notification)) {
        document.body.removeChild(notification);
      }
    }, 3200);
  };

  return (
    <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', overflow: 'auto', zIndex: 1000 }}>
      {/* Development Banner */}
      <div style={{
        position: 'fixed',
        top: '64px',
        left: 0,
        right: 0,
        width: '100%',
        overflow: 'hidden',
        zIndex: 9998,
        background: 'linear-gradient(135deg, hsl(14 97% 49%), hsl(240 76% 8%), hsl(14 97% 45%))',
        color: 'white',
        boxShadow: '0 4px 20px rgba(0, 0, 0, 0.15)'
      }}>
        {/* Animated background pattern */}
        <div style={{
          position: 'absolute',
          inset: 0,
          opacity: 0.15,
          background: 'repeating-linear-gradient(45deg, transparent, transparent 10px, rgba(255,255,255,0.1) 10px, rgba(255,255,255,0.1) 20px)'
        }} />
        
        {/* Animated dots */}
        <div style={{ position: 'absolute', inset: 0, overflow: 'hidden' }}>
          {[...Array(6)].map((_, i) => (
            <div
              key={i}
              style={{
                position: 'absolute',
                width: '8px',
                height: '8px',
                background: 'rgba(255, 255, 255, 0.3)',
                borderRadius: '50%',
                left: `${15 + i * 15}%`,
                top: '50%',
                animation: `bounce 2s infinite ${i * 0.2}s`,
                transform: 'translateY(-50%)'
              }}
            />
          ))}
        </div>

        <div style={{
          position: 'relative',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: windowWidth < 640 ? '6px 16px' : '6px 20px'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            {/* Pulsing warning icon */}
            <div style={{ position: 'relative' }}>
              <AlertTriangle style={{
                width: '16px',
                height: '16px',
                animation: bannerIsAnimating ? 'pulse 1s infinite' : 'none'
              }} />
              <div style={{
                position: 'absolute',
                top: '-2px',
                right: '-2px',
                width: '8px',
                height: '8px',
                background: 'white',
                borderRadius: '50%',
                animation: 'ping 1s infinite',
                opacity: 0.75
              }} />
            </div>

            {/* Main text */}
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '6px',
              fontWeight: '500',
              fontSize: windowWidth < 640 ? '12px' : '13px'
            }}>
              <span style={{ display: windowWidth < 640 ? 'none' : 'inline' }}></span>
              <span style={{
                animation: bannerIsAnimating ? 'pulse 1s infinite' : 'none'
              }}>
                Documentation is Under Active Development
              </span>
            </div>

            {/* Animated icons */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
              <Code style={{
                width: '14px',
                height: '14px',
                animation: 'spin 3s linear infinite'
              }} />
              {/* <Zap style={{
                width: '14px',
                height: '14px',
                animation: bannerIsAnimating ? 'bounce 1s infinite' : 'none'
              }} /> */}
            </div>
          </div>
        </div>

        {/* Bottom accent line */}
        <div style={{
          height: '2px',
          background: 'linear-gradient(to right, transparent, rgba(255,255,255,0.5), transparent)'
        }} />
        
        {/* Shimmer effect */}
        <div style={{
          position: 'absolute',
          inset: 0,
          background: 'linear-gradient(to right, transparent, rgba(255,255,255,0.2), transparent)',
          width: '25%',
          transform: 'translateX(-100%) skewX(-12deg)',
          animation: 'shimmer 3s ease-in-out infinite'
        }} />
      </div>

      <style>{`
        @keyframes shimmer {
          0% { transform: translateX(-100%) skewX(-12deg); }
          100% { transform: translateX(400%) skewX(-12deg); }
        }
        @keyframes bounce {
          0%, 20%, 53%, 80%, 100% { transform: translateY(-50%); }
          40%, 43% { transform: translateY(calc(-50% - 8px)); }
          70% { transform: translateY(calc(-50% - 4px)); }
          90% { transform: translateY(calc(-50% - 2px)); }
        }
        @keyframes pulse {
          0%, 100% { opacity: 1; }
          50% { opacity: 0.5; }
        }
        @keyframes ping {
          75%, 100% { transform: scale(2); opacity: 0; }
        }
        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
        
        /* Custom Scrollbar Styling */
        ::-webkit-scrollbar {
          width: 8px;
        }
        
        ::-webkit-scrollbar-track {
          background: hsl(240 76% 8%);
          border-radius: 4px;
          opacity: 0;
          transition: opacity 0.3s ease;
        }
        
        ::-webkit-scrollbar-thumb {
          background: linear-gradient(135deg, hsl(14 97% 49%), hsl(14 97% 45%));
          border-radius: 4px;
          opacity: 0;
          transition: all 0.3s ease;
        }
        
        ::-webkit-scrollbar-thumb:hover {
          background: linear-gradient(135deg, hsl(14 97% 55%), hsl(14 97% 50%));
          box-shadow: 0 0 10px rgba(251, 146, 60, 0.3);
        }
        
        /* Show scrollbar during scroll */
        body.scrolling ::-webkit-scrollbar-track,
        body.scrolling ::-webkit-scrollbar-thumb {
          opacity: 1;
        }
        
        /* Firefox Scrollbar */
        * {
          scrollbar-width: thin;
          scrollbar-color: transparent transparent;
        }
        
        body.scrolling * {
          scrollbar-color: hsl(14 97% 49%) hsl(240 76% 8%);
        }
        
        /* Hide external link icons if needed */
        a[hideexternalicon] {
          /* Styles for links with hideexternalicon attribute */
        }
        
        /* Get Started Button Animations - Same as Banner */
        .get-started-btn .btn-dot {
          animation: none;
        }
        
        .get-started-btn:hover .btn-dot {
          animation: bounce 2s infinite;
        }
        
        .get-started-btn:hover .btn-dot:nth-child(1) {
          animation-delay: 0s;
        }
        
        .get-started-btn:hover .btn-dot:nth-child(2) {
          animation-delay: 0.2s;
        }
        
        .get-started-btn:hover .btn-dot:nth-child(3) {
          animation-delay: 0.4s;
        }
        
        .get-started-btn .btn-shimmer {
          animation: none;
        }
        
        .get-started-btn:hover .btn-shimmer {
          animation: shimmer 3s ease-in-out infinite;
        }
        
        .get-started-btn .btn-icon {
          animation: none;
        }
        
        .get-started-btn:hover .btn-icon {
          animation: spin 3s linear infinite;
        }
        
        .get-started-btn:hover {
          background: linear-gradient(135deg, hsl(14 97% 55%), hsl(14 97% 49%));
          box-shadow: 0 8px 25px rgba(251, 146, 60, 0.4);
        }
        
        .get-started-btn {
          transition: background 0.2s ease, box-shadow 0.2s ease !important;
        }
        
        /* Prevent iframe flashing */
        iframe {
          display: block;
          background: #03032a !important;
        }
        
        iframe[src*="gopher-animation"] {
          background: #03032a !important;
          opacity: 1 !important;
        }
        
        /* Immediately hide Vocs elements to prevent flash */
        .vocs_Header,
        .vocs_Nav,
        .vocs_TopNav {
          display: none !important;
        }
        
        /* Immediately style containers */
        [data-layout="landing"],
        main,
        article,
        .vocs_Content,
        .vocs_ContentWrapper {
          padding: 0 !important;
          margin: 0 !important;
          max-width: none !important;
          width: 100% !important;
        }
        
        .bg-background { background-color: hsl(240 76% 8%); }
        .bg-primary { background-color: hsl(14 97% 49%); }
        .text-primary-foreground { color: hsl(210 40% 98%); }
        .border-gray-700 { border-color: rgb(55 65 81); }
        .text-gray-400 { color: rgb(156 163 175); }
        .text-white { color: rgb(255 255 255); }
        .text-gray-900 { color: rgb(17 24 39); }
        .text-gray-300 { color: rgb(209 213 219); }
        .text-green-400 { color: rgb(74 222 128); }
        .bg-white { background-color: rgb(255 255 255); }
        .bg-gray-800 { background-color: rgb(31 41 55); }
        .hover\\:bg-primary\\/90:hover { background-color: hsl(14 97% 49% / 0.9); }
        .hover\\:text-white:hover { color: rgb(255 255 255); }
        .hover\\:bg-white\\/90:hover { background-color: rgb(255 255 255 / 0.9); }
        .hover\\:bg-gray-50:hover { background-color: rgb(249 250 251); }
        .hover\\:bg-gray-100:hover { background-color: rgb(243 244 246); }
        .transition-colors { transition-property: color, background-color, border-color, text-decoration-color, fill, stroke; transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1); transition-duration: 150ms; }
        .backdrop-blur-sm { backdrop-filter: blur(4px); }
        .border-b { border-bottom-width: 1px; }
        .border-gray-200 { border-color: rgb(229 231 235); }
        .border-gray-300 { border-color: rgb(209 213 219); }
        .border { border-width: 1px; }
        .rounded-md { border-radius: 0.375rem; }
        .rounded-lg { border-radius: 0.5rem; }
        .rounded-full { border-radius: 9999px; }
        .px-4 { padding-left: 1rem; padding-right: 1rem; }
        .py-2 { padding-top: 0.5rem; padding-bottom: 0.5rem; }
        .px-8 { padding-left: 2rem; padding-right: 2rem; }
        .py-3 { padding-top: 0.75rem; padding-bottom: 0.75rem; }
        .px-3 { padding-left: 0.75rem; padding-right: 0.75rem; }
        .py-16 { padding-top: 4rem; padding-bottom: 4rem; }
        .pt-8 { padding-top: 2rem; }
        .text-sm { font-size: 0.875rem; line-height: 1.25rem; }
        .text-lg { font-size: 1.125rem; line-height: 1.75rem; }
        .text-xl { font-size: 1.25rem; line-height: 1.75rem; }
        .text-xs { font-size: 0.75rem; line-height: 1rem; }
        .font-medium { font-weight: 500; }
        .font-bold { font-weight: 700; }
        .font-semibold { font-weight: 600; }
        .font-mono { font-family: ui-monospace, SFMono-Regular, "SF Mono", Monaco, Inconsolata, "Roboto Mono", monospace; }
        .leading-relaxed { line-height: 1.625; }
        .shadow-2xl { box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.25); }
        .animate-pulse { animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite; }
        .bg-gradient-to-br { background-image: linear-gradient(to bottom right, var(--tw-gradient-stops)); }
        .from-blue-400 { --tw-gradient-from: #60a5fa var(--tw-gradient-from-position); --tw-gradient-to: rgb(96 165 250 / 0) var(--tw-gradient-to-position); --tw-gradient-stops: var(--tw-gradient-from), var(--tw-gradient-to); }
        .to-blue-600 { --tw-gradient-to: #2563eb var(--tw-gradient-to-position); }
        .from-pink-300 { --tw-gradient-from: #f9a8d4 var(--tw-gradient-from-position); --tw-gradient-to: rgb(249 168 212 / 0) var(--tw-gradient-to-position); --tw-gradient-stops: var(--tw-gradient-from), var(--tw-gradient-to); }
        .to-white { --tw-gradient-to: #fff var(--tw-gradient-to-position); }
        .bg-gradient-to-r { background-image: linear-gradient(to right, var(--tw-gradient-stops)); }
        .from-primary { --tw-gradient-from: hsl(14 97% 49%) var(--tw-gradient-from-position); --tw-gradient-to: hsl(14 97% 49% / 0) var(--tw-gradient-to-position); --tw-gradient-stops: var(--tw-gradient-from), var(--tw-gradient-to); }
        .to-primary\\/80 { --tw-gradient-to: hsl(14 97% 49% / 0.8) var(--tw-gradient-to-position); }
        .bg-clip-text { background-clip: text; }
        .text-transparent { color: transparent; }
        
        @keyframes pulse {
          50% { opacity: .5; }
        }
      `}</style>

      {/* Header */}
      <header id="starknet-header" style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        zIndex: 9999,
        backgroundColor: 'rgb(255 255 255 / 0.95)',
        backdropFilter: 'blur(4px)',
        borderBottom: '1px solid rgb(229 231 235)'
      }}>
        <div style={{ maxWidth: '80rem', margin: '0 auto', padding: '0 1rem' }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', height: '4rem' }}>
            {/* Logo */}
            <div style={{ display: 'flex', alignItems: 'center' }}>
              <img
                src="/Starknet.Go_Horizontal_Light.svg"
                alt="Starknet.go Logo"
                style={{
                  height: '36px',
                  objectFit: 'contain'
                }}
              />
            </div>

            {/* Navigation */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <a 
                href="/docs/introduction/why-starknet-go" 
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: '0.25rem',
                  padding: '0.375rem 0.75rem',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  borderRadius: '0.5rem',
                  border: '1px solid hsl(240 76% 8%)',
                  backgroundColor: 'white',
                  color: 'rgb(17 24 39)',
                  cursor: 'pointer',
                  transition: 'all 0.15s',
                  textDecoration: 'none'
                }}
                onMouseEnter={(e) => (e.target as HTMLElement).style.backgroundColor = 'rgb(249 250 251)'}
                onMouseLeave={(e) => (e.target as HTMLElement).style.backgroundColor = 'white'}
              >
                <img src="/docIcon.png" alt="" style={{ width: '14px', height: '14px' }} />
                DOCS
              </a>
              
              {/* Twitter/X */}
              <a href="https://x.com/NethermindStark" target="_blank" rel="noopener noreferrer"
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: '2rem',
                  height: '2rem',
                  borderRadius: '0.5rem',
                  color: 'rgb(17 24 39)',
                  cursor: 'pointer',
                  transition: 'all 0.15s',
                  textDecoration: 'none'
                }}
                onMouseEnter={(e) => (e.target as HTMLElement).style.backgroundColor = 'rgb(243 244 246)'}
                onMouseLeave={(e) => (e.target as HTMLElement).style.backgroundColor = 'transparent'}>
                <img src="/xicon.svg" alt="X" style={{ width: '16px', height: '16px' }} />
              </a>
              
              {/* Telegram */}
              <a href="https://t.me/StarknetGo" target="_blank" rel="noopener noreferrer"
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: '2rem',
                  height: '2rem',
                  borderRadius: '0.5rem',
                  color: 'rgb(17 24 39)',
                  cursor: 'pointer',
                  transition: 'all 0.15s',
                  textDecoration: 'none'
                }}
                onMouseEnter={(e) => (e.target as HTMLElement).style.backgroundColor = 'rgb(243 244 246)'}
                onMouseLeave={(e) => (e.target as HTMLElement).style.backgroundColor = 'transparent'}>
                <img src="/Telegram.svg" alt="Telegram" style={{ width: '16px', height: '16px' }} />
              </a>
              
              {/* Github */}
              <a href="https://github.com/NethermindEth/starknet.go" target="_blank" rel="noopener noreferrer"
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: '2rem',
                  height: '2rem',
                  borderRadius: '0.5rem',
                  color: 'rgb(17 24 39)',
                  cursor: 'pointer',
                  transition: 'all 0.15s'
                }}
                onMouseEnter={(e) => (e.target as HTMLElement).style.backgroundColor = 'rgb(243 244 246)'}
                onMouseLeave={(e) => (e.target as HTMLElement).style.backgroundColor = 'transparent'}>
                <img src="/Githubicon.svg" alt="GitHub" style={{ width: '16px', height: '16px' }} />
              </a>
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        position: 'relative',
        overflow: 'hidden',
        paddingTop: '114px'
      }}>
        {/* Main background */}
        <div style={{ position: 'absolute', inset: 0, backgroundColor: '#03032a' }} />
        
        {/* Enhanced gradient elements - positioned away from gopher animation */}
        <div style={{ position: 'absolute', inset: 0 }}>
          {/* Top left gradient */}
          <div style={{
            position: 'absolute',
            width: '60rem',
            height: '60rem',
            borderRadius: '50%',
            opacity: 0.25,
            background: 'radial-gradient(circle, rgba(255, 255, 255, 0.12) 0%, transparent 65%)',
            top: '-30rem',
            left: '-25rem',
            filter: 'blur(50px)'
          }} />
          
          {/* Bottom right gradient */}
          <div style={{
            position: 'absolute',
            width: '55rem',
            height: '55rem',
            borderRadius: '50%',
            opacity: 0.2,
            background: 'radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 65%)',
            bottom: '-25rem',
            right: '-20rem',
            filter: 'blur(45px)'
          }} />
          
          {/* Bottom left accent gradient */}
          <div style={{
            position: 'absolute',
            width: '40rem',
            height: '40rem',
            borderRadius: '50%',
            opacity: 0.15,
            background: 'radial-gradient(circle, rgba(251, 146, 60, 0.15) 0%, transparent 65%)',
            bottom: '-15rem',
            left: '-10rem',
            filter: 'blur(35px)'
          }} />
          
          {/* Top right accent */}
          <div style={{
            position: 'absolute',
            width: '35rem',
            height: '35rem',
            borderRadius: '50%',
            opacity: 0.1,
            background: 'radial-gradient(circle, rgba(251, 146, 60, 0.12) 0%, transparent 65%)',
            top: '-10rem',
            right: '-12rem',
            filter: 'blur(30px)'
          }} />
          
          {/* Subtle mesh gradient overlay - avoiding center area */}
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: `
              radial-gradient(ellipse 70% 50% at 10% 15%, rgba(255, 255, 255, 0.04) 0%, transparent 70%),
              radial-gradient(ellipse 60% 40% at 90% 85%, rgba(251, 146, 60, 0.03) 0%, transparent 70%),
              radial-gradient(ellipse 50% 30% at 5% 95%, rgba(255, 255, 255, 0.03) 0%, transparent 70%)
            `,
            filter: 'blur(30px)'
          }} />
        </div>
        
        <div style={{
          position: 'relative',
          zIndex: 10,
          textAlign: 'center',
          padding: windowWidth < 640 ? '0 1rem' : '0 1.5rem',
          maxWidth: '64rem',
          margin: '0 auto'
        }}>
          {/* Gopher Animation */}
          <div style={{ 
            marginBottom: '2rem', 
            display: 'flex', 
            justifyContent: 'center',
            alignItems: 'center'
          }}>
            <iframe
              src="/gopher-animation.html"
              style={{
                width: '300px',
                height: '195px',
                border: 'none',
                background: '#03032a'
              }}
              title="Gopher Stars Animation"
            />
          </div>

          {/* Release Badge */}
          <div style={{ marginBottom: '1.5rem', display: 'flex', justifyContent: 'center' }}>
            <div style={{
              backgroundColor: 'white',
              color: 'rgb(17 24 39)',
              padding: '0.375rem 0.875rem',
              borderRadius: '9999px',
              display: 'flex',
              alignItems: 'center',
              gap: '0.375rem',
              transition: 'all 0.15s',
              cursor: 'pointer'
            }}
            onMouseEnter={(e) => (e.target as HTMLElement).style.backgroundColor = 'rgb(255 255 255 / 0.9)'}
            onMouseLeave={(e) => (e.target as HTMLElement).style.backgroundColor = 'white'}>
              <span style={{ fontSize: '0.75rem', fontWeight: '600' }}>LATEST</span>
              <span style={{ fontSize: '0.75rem' }}>Starknet.go v0.10.0 Released</span>
              <ArrowRight style={{ width: '0.625rem', height: '0.625rem' }} />
            </div>
          </div>

          {/* Main Heading */}
          <h1 style={{
            fontSize: windowWidth < 640 ? '2.5rem' : windowWidth < 1024 ? '3.5rem' : '4rem',
            fontWeight: '700',
            marginBottom: '1.25rem',
            background: 'linear-gradient(to right, hsl(14 97% 49%), hsl(14 97% 49% / 0.8))',
            backgroundClip: 'text',
            color: 'transparent',
            lineHeight: '1.2',
            letterSpacing: '-0.025em'
          }}>
            Starknet for Go Developers
          </h1>

          {/* Subtitle */}
          <p style={{
            fontSize: windowWidth < 640 ? '1rem' : '1.125rem',
            lineHeight: '1.5',
            color: 'rgb(209 213 219)',
            marginBottom: '2.5rem',
            maxWidth: '36rem',
            margin: '0 auto 2.5rem auto'
          }}>
            Building the future of Starknet with a powerful Go implementation<br />
            for scalable and efficient blockchain development.
          </p>

          {/* Code Snippet */}
          <div style={{ marginBottom: '1.5rem', maxWidth: '28rem', margin: '0 auto 1.5rem auto' }}>
            <div style={{
              backgroundColor: 'rgb(31 41 55)',
              border: '1px solid rgb(55 65 81)',
              borderRadius: '0.5rem',
              padding: '0.875rem 1rem',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between'
            }}>
              <code style={{
                color: 'rgb(74 222 128)',
                fontFamily: 'ui-monospace, SFMono-Regular, "SF Mono", Monaco, Inconsolata, "Roboto Mono", monospace',
                fontSize: windowWidth < 640 ? '0.75rem' : '0.875rem'
              }}>
                go get github.com/NethermindEth/starknet.go
              </code>
              <button
                onClick={copyToClipboard}
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: '2rem',
                  height: '2rem',
                  borderRadius: '0.375rem',
                  backgroundColor: 'transparent',
                  border: 'none',
                  color: 'rgb(156 163 175)',
                  cursor: 'pointer',
                  transition: 'all 0.15s',
                  marginLeft: '0.5rem'
                }}
                onMouseEnter={(e) => (e.target as HTMLElement).style.color = 'white'}
                onMouseLeave={(e) => (e.target as HTMLElement).style.color = 'rgb(156 163 175)'}
              >
                <Copy style={{ width: '0.875rem', height: '0.875rem' }} />
              </button>
            </div>
          </div>

          {/* CTA Button */}
          <a 
            href="/docs/introduction/getting-started"
            className="get-started-btn"
            style={{
              display: 'inline-flex',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '0.5rem',
              backgroundColor: 'hsl(14 97% 49%)',
              color: 'white',
              padding: '0.75rem 2rem',
              fontSize: '1.125rem',
              fontWeight: '600',
              borderRadius: '0.5rem',
              border: 'none',
              cursor: 'pointer',
              textDecoration: 'none',
              position: 'relative',
              overflow: 'hidden',
              boxShadow: '0 4px 15px rgba(251, 146, 60, 0.2)'
            }}>
            {/* Animated background pattern */}
            <div style={{
              position: 'absolute',
              inset: 0,
              opacity: 0.15,
              background: 'repeating-linear-gradient(45deg, transparent, transparent 10px, rgba(255,255,255,0.1) 10px, rgba(255,255,255,0.1) 20px)'
            }} />
            
            {/* Animated dots */}
            <div style={{ position: 'absolute', inset: 0, overflow: 'hidden' }}>
              {[...Array(3)].map((_, i) => (
                <div
                  key={i}
                  className="btn-dot"
                  style={{
                    position: 'absolute',
                    width: '4px',
                    height: '4px',
                    background: 'rgba(255, 255, 255, 0.4)',
                    borderRadius: '50%',
                    left: `${20 + i * 25}%`,
                    top: '50%',
                    transform: 'translateY(-50%)'
                  }}
                />
              ))}
            </div>

            {/* Shimmer effect */}
            <div className="btn-shimmer" style={{
              position: 'absolute',
              inset: 0,
              background: 'linear-gradient(to right, transparent, rgba(255,255,255,0.3), transparent)',
              width: '25%',
              transform: 'translateX(-100%) skewX(-12deg)'
            }} />
            
            <span style={{ position: 'relative', zIndex: 1 }}>Get Started</span>
            <Code style={{
              width: '16px',
              height: '16px',
              position: 'relative',
              zIndex: 1
            }} className="btn-icon" />
          </a>
        </div>
      </section>

      {/* Footer */}
      <footer style={{ backgroundColor: 'hsl(240 76% 8%)', borderTop: '1px solid rgb(55 65 81)' }}>
        <div style={{ maxWidth: '80rem', margin: '0 auto', padding: '4rem 1rem' }}>
          <div style={{
            display: 'flex',
            flexDirection: windowWidth < 768 ? 'column' : 'row',
            justifyContent: 'space-evenly',
            alignItems: 'flex-start',
            gap: windowWidth < 768 ? '2rem' : '1rem',
            marginBottom: '2rem'
          }}>
            {/* Solutions */}
            <div>
              <h3 style={{ color: 'white', fontWeight: '600', marginBottom: '1rem' }}>Solutions</h3>
              <ul style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                {[
                  { name: "Starknet.go", href: "https://starknet-go.nethermind.io/#" },
                  { name: "Voyager", href: "https://voyager.online/" },
                  { name: "Juno", href: "https://www.nethermind.io/juno" },
                  { name: "CairoVM - Go", href: "https://github.com/NethermindEth/cairo-vm-go" },
                  { name: "Starkweb", href: "https://www.starkweb.xyz/" }
                ].map((item) => (
                  <li key={item.name}>
                    <a
                      href={item.href}
                      style={{
                        color: 'rgb(156 163 175)',
                        textDecoration: 'none',
                        transition: 'color 0.15s'
                      }}
                      onMouseEnter={(e) => (e.target as HTMLElement).style.color = 'white'}
                      onMouseLeave={(e) => (e.target as HTMLElement).style.color = 'rgb(156 163 175)'}
                    >
                      {item.name}
                    </a>
                  </li>
                ))}
              </ul>
            </div>

            {/* Company */}
            <div>
              <h3 style={{ color: 'white', fontWeight: '600', marginBottom: '1rem' }}>Company</h3>
              <ul style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                {[
                  { name: "About Nethermind", href: "https://www.nethermind.io/" },
                  { name: "Blog", href: "https://www.nethermind.io/blog" },
                  { name: "Careers", href: "https://www.nethermind.io/open-roles" },
                  { name: "Events", href: "https://www.nethermind.io/events" }
                ].map((item) => (
                  <li key={item.name}>
                    <a
                      href={item.href}
                      style={{
                        color: 'rgb(156 163 175)',
                        textDecoration: 'none',
                        transition: 'color 0.15s'
                      }}
                      onMouseEnter={(e) => (e.target as HTMLElement).style.color = 'white'}
                      onMouseLeave={(e) => (e.target as HTMLElement).style.color = 'rgb(156 163 175)'}
                    >
                      {item.name}
                    </a>
                  </li>
                ))}
              </ul>
            </div>


            {/* Ecosystem */}
            <div>
              <h3 style={{ color: 'white', fontWeight: '600', marginBottom: '1rem' }}>Ecosystem</h3>
              <ul style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                {[
                  { name: "Starknet", href: "https://www.starknet.io/" },
                  { name: "Github", href: "https://github.com/NethermindEth/starknet.go" },
                  { name: "X", href: "https://x.com/NethermindStark?ref_src=twsrc%5Egoogle%7Ctwcamp%5Eserp%7Ctwgr%5Eauthor" },
                  { name: "Telegram", href: "https://t.me/StarknetGo" }
                ].map((item) => (
                  <li key={item.name}>
                    <a
                      href={item.href}
                      style={{
                        color: 'rgb(156 163 175)',
                        textDecoration: 'none',
                        transition: 'color 0.15s'
                      }}
                      onMouseEnter={(e) => (e.target as HTMLElement).style.color = 'white'}
                      onMouseLeave={(e) => (e.target as HTMLElement).style.color = 'rgb(156 163 175)'}
                    >
                      {item.name}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Bottom Section */}
          <div style={{
            marginTop: '3rem',
            paddingTop: '2rem',
            borderTop: '1px solid rgb(55 65 81)',
            display: 'flex',
            flexDirection: windowWidth < 640 ? 'column' : 'row',
            justifyContent: 'space-between',
            alignItems: 'center'
          }}>
            <img
              src="/Starknet.Go_Dark_Powered_by_Nethermind.png"
              alt="Starknet.go powered by Nethermind"
              style={{ height: '2rem' }}
            />
            
            <span style={{ 
              color: 'rgb(156 163 175)', 
              fontSize: '0.875rem',
              marginTop: windowWidth < 640 ? '1rem' : 0
            }}>
              © 2025 Nethermind. All Rights Reserved.
            </span>
          </div>
        </div>
      </footer>
    </div>
  );
}