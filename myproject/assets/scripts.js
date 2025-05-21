document.querySelectorAll('.tab-button').forEach(button => {
    button.addEventListener('click', () => {
      // Удалить активный класс у всех кнопок
      document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
  
      // Скрыть все вкладки
      document.querySelectorAll('.tab-content').forEach(tab => tab.classList.remove('active'));
  
      // Активировать текущую вкладку и кнопку
      const tabId = button.dataset.tab;
      button.classList.add('active');
      document.getElementById(tabId).classList.add('active');
    });
  });
  
  document.querySelectorAll('.accordion-header').forEach(header => {
    header.addEventListener('click', () => {
      const content = header.nextElementSibling;
  
      // Развернуть или свернуть контент
      if (content.style.display === 'block') {
        content.style.display = 'none';
      } else {
        content.style.display = 'block';
      }
    });
  });
  