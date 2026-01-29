export interface Lesson {
  id: number;
  title: string;
  description: string;
  theory: string;
  example: string;
  starterCode: string;
  challenge: string;
  hints: string[];
  validation: {
    type: 'contains' | 'exact' | 'regex';
    value: string | RegExp;
    errorMessage?: string;
  };
  xpReward: number;
}

export const lessons: Lesson[] = [
  {
    id: 1,
    title: 'Your First Tag: Paragraphs',
    description: 'Learn the most basic HTML tag - the paragraph!',
    theory: 'The `<p>` tag creates a paragraph of text. It\'s one of the most common HTML tags you\'ll use. Everything between `<p>` and `</p>` becomes a paragraph.',
    example: '<p>Hello, World!</p>',
    starterCode: '',
    challenge: 'Create a paragraph that says "Hello, HTML!"',
    hints: [
      'Start with an opening tag: <p>',
      'Type your text in the middle',
      'End with a closing tag: </p>',
    ],
    validation: {
      type: 'contains',
      value: '<p>',
      errorMessage: 'Make sure to use the <p> tag!',
    },
    xpReward: 25,
  },
  {
    id: 2,
    title: 'Headings: Making Titles',
    description: 'Learn how to create headings from big to small',
    theory: 'HTML has six heading levels: `<h1>` through `<h6>`. `<h1>` is the biggest and most important, `<h6>` is the smallest. Use headings to structure your content!',
    example: '<h1>Big Title</h1>\n<h2>Smaller Title</h2>\n<h3>Even Smaller</h3>',
    starterCode: '',
    challenge: 'Create a big heading (h1) that says "Welcome to HTML"',
    hints: [
      'Use the <h1> tag for the biggest heading',
      'Don\'t forget to close it with </h1>',
      'Headings can contain any text',
    ],
    validation: {
      type: 'contains',
      value: '<h1>',
      errorMessage: 'Use the <h1> tag for your heading!',
    },
    xpReward: 25,
  },
  {
    id: 3,
    title: 'Bold and Italic: Emphasizing Text',
    description: 'Make text stand out with bold and italic',
    theory: 'Use `<strong>` to make text **bold** and `<em>` to make text *italic*. These tags help emphasize important words!',
    example: '<p>This is <strong>bold</strong> and this is <em>italic</em>!</p>',
    starterCode: '<p></p>',
    challenge: 'Make a paragraph with some bold text using <strong>',
    hints: [
      'Put <strong> tags around the text you want bold',
      'Example: <strong>important</strong>',
      'Don\'t forget the closing </strong> tag',
    ],
    validation: {
      type: 'contains',
      value: '<strong>',
      errorMessage: 'Use <strong> to make text bold!',
    },
    xpReward: 25,
  },
  {
    id: 4,
    title: 'Lists: Organizing Information',
    description: 'Create bullet points and numbered lists',
    theory: 'Lists help organize information. Use `<ul>` for bullet points (unordered list) and `<ol>` for numbered lists (ordered list). Each item goes in an `<li>` tag.',
    example: '<ul>\n  <li>First item</li>\n  <li>Second item</li>\n  <li>Third item</li>\n</ul>',
    starterCode: '',
    challenge: 'Create an unordered list with at least 2 items',
    hints: [
      'Start with <ul> for an unordered list',
      'Each item needs <li> tags',
      'Close with </ul> at the end',
    ],
    validation: {
      type: 'contains',
      value: '<ul>',
      errorMessage: 'Create an unordered list with <ul>!',
    },
    xpReward: 25,
  },
  {
    id: 5,
    title: 'Links: Connecting Pages',
    description: 'Learn how to create clickable links',
    theory: 'The `<a>` tag creates links. Use the `href` attribute to specify where the link goes. Links are what make the web... a web!',
    example: '<a href="https://example.com">Click me!</a>',
    starterCode: '',
    challenge: 'Create a link that says "Visit Google" and links to https://google.com',
    hints: [
      'Use the <a> tag with an href attribute',
      'Format: <a href="url">link text</a>',
      'Remember to include https:// in the URL',
    ],
    validation: {
      type: 'contains',
      value: '<a href=',
      errorMessage: 'Create a link using <a href="...">!',
    },
    xpReward: 25,
  },
  {
    id: 6,
    title: 'Images: Adding Pictures',
    description: 'Display images on your webpage',
    theory: 'The `<img>` tag displays images. It needs a `src` attribute (the image URL) and an `alt` attribute (description for accessibility). Note: `<img>` doesn\'t have a closing tag!',
    example: '<img src="https://via.placeholder.com/150" alt="A placeholder image">',
    starterCode: '',
    challenge: 'Add an image with alt text describing it',
    hints: [
      'Use <img src="..." alt="...">',
      'No closing tag needed for img!',
      'Always include alt text for accessibility',
    ],
    validation: {
      type: 'contains',
      value: '<img',
      errorMessage: 'Add an image using the <img> tag!',
    },
    xpReward: 25,
  },
  {
    id: 7,
    title: 'Line Breaks: Spacing Things Out',
    description: 'Control spacing with breaks and horizontal rules',
    theory: 'Use `<br>` for line breaks (like pressing Enter) and `<hr>` for horizontal lines. These are self-closing tags - no closing tag needed!',
    example: '<p>Line one<br>Line two</p>\n<hr>\n<p>After the line</p>',
    starterCode: '<p>First line</p>\n<p>Second line</p>',
    challenge: 'Add a line break between two lines of text',
    hints: [
      'Use <br> to create a line break',
      '<br> doesn\'t need a closing tag',
      'Place it where you want the break',
    ],
    validation: {
      type: 'contains',
      value: '<br>',
      errorMessage: 'Use <br> for a line break!',
    },
    xpReward: 25,
  },
  {
    id: 8,
    title: 'Divs: Grouping Content',
    description: 'Organize your page with containers',
    theory: 'The `<div>` tag is a container for grouping content. It\'s like a box that holds other elements. Divs help structure your page layout.',
    example: '<div>\n  <h2>Section Title</h2>\n  <p>Section content</p>\n</div>',
    starterCode: '',
    challenge: 'Create a div that contains a heading and a paragraph',
    hints: [
      'Start with <div>',
      'Put other tags inside it',
      'Close with </div>',
    ],
    validation: {
      type: 'contains',
      value: '<div>',
      errorMessage: 'Create a container with <div>!',
    },
    xpReward: 25,
  },
  {
    id: 9,
    title: 'Spans: Inline Styling',
    description: 'Target specific words or phrases',
    theory: 'The `<span>` tag is like a mini-container for inline content. Use it when you want to style or target specific words within a paragraph.',
    example: '<p>This is <span style="color: red;">red text</span> in a paragraph.</p>',
    starterCode: '<p></p>',
    challenge: 'Create a paragraph with a span around some words',
    hints: [
      'Put <span> tags around specific words',
      'Span works inside other tags like <p>',
      'Example: <span>highlighted</span>',
    ],
    validation: {
      type: 'contains',
      value: '<span>',
      errorMessage: 'Use <span> to wrap some text!',
    },
    xpReward: 25,
  },
  {
    id: 10,
    title: 'Basic Structure: HTML Document',
    description: 'Learn the skeleton of every HTML page',
    theory: 'Every HTML page has three main parts: `<html>` wraps everything, `<head>` contains meta information, and `<body>` contains visible content.',
    example: '<html>\n  <head>\n    <title>My Page</title>\n  </head>\n  <body>\n    <h1>Content goes here</h1>\n  </body>\n</html>',
    starterCode: '',
    challenge: 'Create a basic HTML structure with html, head, and body tags',
    hints: [
      'Start with <html>',
      'Add <head> and </head>',
      'Add <body> and </body>',
      'Close with </html>',
    ],
    validation: {
      type: 'contains',
      value: '<html>',
      errorMessage: 'Start with the <html> tag!',
    },
    xpReward: 25,
  },
  {
    id: 11,
    title: 'Meta Tags: Page Information',
    description: 'Add metadata to your HTML page',
    theory: 'Meta tags go in the `<head>` and provide information about your page. The `<title>` tag sets the page title (shown in browser tabs). Meta tags don\'t have closing tags.',
    example: '<head>\n  <title>My Awesome Page</title>\n  <meta charset="UTF-8">\n</head>',
    starterCode: '<head>\n</head>',
    challenge: 'Add a title tag to the head section',
    hints: [
      'Use <title> inside <head>',
      'The title appears in the browser tab',
      'Format: <title>Your Title</title>',
    ],
    validation: {
      type: 'contains',
      value: '<title>',
      errorMessage: 'Add a <title> tag in the head!',
    },
    xpReward: 25,
  },
  {
    id: 12,
    title: 'Forms Part 1: Input and Buttons',
    description: 'Create interactive forms with inputs',
    theory: 'Forms collect user input. Use `<form>` to wrap form elements, `<input>` for text fields, and `<button>` for clickable buttons.',
    example: '<form>\n  <input type="text" placeholder="Enter name">\n  <button>Submit</button>\n</form>',
    starterCode: '<form>\n</form>',
    challenge: 'Add an input field and a button to the form',
    hints: [
      'Use <input type="text">',
      'Add <button>Button Text</button>',
      'Both go inside the <form> tags',
    ],
    validation: {
      type: 'contains',
      value: '<input',
      errorMessage: 'Add an <input> field to your form!',
    },
    xpReward: 25,
  },
  {
    id: 13,
    title: 'Forms Part 2: Textarea and Select',
    description: 'Add more form elements',
    theory: '`<textarea>` creates multi-line text fields. `<select>` creates dropdown menus with `<option>` tags for each choice. `<label>` describes form fields.',
    example: '<form>\n  <label>Comments:</label>\n  <textarea rows="3"></textarea>\n  <select>\n    <option>Choice 1</option>\n    <option>Choice 2</option>\n  </select>\n</form>',
    starterCode: '<form>\n</form>',
    challenge: 'Add a textarea or a select dropdown to your form',
    hints: [
      'For textarea: <textarea></textarea>',
      'For dropdown: <select><option>Item</option></select>',
      'Add rows="3" to textarea for height',
    ],
    validation: {
      type: 'contains',
      value: '<textarea',
      errorMessage: 'Add a <textarea> or <select> to your form!',
    },
    xpReward: 25,
  },
  {
    id: 14,
    title: 'Tables: Organizing Data',
    description: 'Create tables with rows and columns',
    theory: 'Tables display data in rows and columns. Use `<table>` to start, `<tr>` for rows, `<th>` for headers, and `<td>` for data cells.',
    example: '<table>\n  <tr>\n    <th>Name</th>\n    <th>Age</th>\n  </tr>\n  <tr>\n    <td>John</td>\n    <td>25</td>\n  </tr>\n</table>',
    starterCode: '<table>\n</table>',
    challenge: 'Create a table with at least one row of data',
    hints: [
      'Use <tr> for each row',
      'Use <td> for data cells',
      'Example: <tr><td>Data</td></tr>',
    ],
    validation: {
      type: 'contains',
      value: '<tr>',
      errorMessage: 'Add table rows with <tr>!',
    },
    xpReward: 25,
  },
  {
    id: 15,
    title: 'Final Challenge: Complete Webpage',
    description: 'Build a complete HTML page with everything you\'ve learned!',
    theory: 'Time to put it all together! Create a complete webpage using all the tags you\'ve learned: structure, headings, paragraphs, links, images, lists, and more.',
    example: '<!DOCTYPE html>\n<html>\n<head>\n  <title>My Page</title>\n</head>\n<body>\n  <h1>Welcome</h1>\n  <p>This is my first webpage!</p>\n  <ul>\n    <li>Item 1</li>\n    <li>Item 2</li>\n  </ul>\n</body>\n</html>',
    starterCode: '',
    challenge: 'Create a complete webpage with: html structure, a title, a heading, a paragraph, and a list',
    hints: [
      'Start with <!DOCTYPE html>',
      'Include <html>, <head>, <body>',
      'Add a <title> in the head',
      'Add content: <h1>, <p>, <ul> in the body',
    ],
    validation: {
      type: 'contains',
      value: '<html>',
      errorMessage: 'Create a complete HTML structure!',
    },
    xpReward: 50,
  },
];

export function getLessonById(id: number): Lesson | undefined {
  return lessons.find(lesson => lesson.id === id);
}

export function validateLessonCode(lessonId: number, code: string): boolean {
  const lesson = getLessonById(lessonId);
  if (!lesson) return false;

  const validation = lesson.validation;

  switch (validation.type) {
    case 'contains':
      return code.toLowerCase().includes(validation.value.toString().toLowerCase());

    case 'exact':
      return code.trim() === validation.value.toString().trim();

    case 'regex':
      if (validation.value instanceof RegExp) {
        return validation.value.test(code);
      }
      return false;

    default:
      return false;
  }
}
