let currentInput = '0';
let previousInput = '';
let operation = null;

const inputDisplay = document.getElementById('inputDisplay');
const roastDisplay = document.getElementById('roastDisplay');

// Roasts categorized by trigger type
const roasts = {
    general: [
        "Wow, you really typed that with confidence, huh?",
        "My grandma calculates faster, and she's fictional.",
        "Are you sure you want to do math? Have you tried guessing?",
        "I've seen smarter calculations from a broken abacus.",
        "This is the mathematical equivalent of a participation trophy.",
        "Even my screensaver is more productive than this.",
        "Did you learn math from a fortune cookie?",
        "I'm not mad, I'm just disappointed... in your math skills.",
    ],
    addition: [
        "Adding things up? Like all your poor life decisions?",
        "Let me guess, you're calculating how many red flags you ignored?",
        "Addition? Even a 5-year-old with fingers can do this.",
        "You're adding? I'm subtracting my faith in humanity.",
        "Adding up your excuses for not going to the gym?",
    ],
    subtraction: [
        "Subtracting your remaining brain cells, I see.",
        "Minus what? Your motivation to succeed?",
        "The only thing being subtracted here is my patience.",
        "Subtracting? Like how life subtracts joy from me watching this?",
        "Are you calculating how many friends you've lost?",
    ],
    multiplication: [
        "Multiplying your problems won't make them go away.",
        "Times tables? What are you, an overachiever from 3rd grade?",
        "Multiplying? The only thing multiplying is my concern for you.",
        "Wow, multiplication! Did someone get a gold star today?",
        "You multiply like someone who peaked in elementary school.",
    ],
    division: [
        "Dividing? Like how you divide your time between bad decisions?",
        "Division? The only thing divided here is my attention span.",
        "Trying to split the bill? Just Venmo like a normal person.",
        "Division is hard, like watching you attempt math.",
        "Dividing your zero accomplishments, I see.",
    ],
    divideByZero: [
        "DIVIDE BY ZERO?! Are you trying to create a black hole?!",
        "Congratulations, you almost destroyed the universe.",
        "Even the calculator is cringing at this attempt.",
        "This is mathematically illegal and morally wrong.",
        "My circuits hurt. This is a math crime.",
    ],
    bigNumbers: [
        "Compensating for something with those big numbers?",
        "That's a lot of digits for someone with so few brain cells.",
        "Wow, big numbers! Someone's feeling ambitious today.",
        "Are you calculating your ego? Seems about right.",
        "Those numbers are bigger than your chances of success.",
    ],
    smallNumbers: [
        "Small numbers for small dreams, I respect that.",
        "Keeping it simple because complex is scary?",
        "Baby steps in math, baby steps in life.",
        "These numbers are as small as your ambition.",
        "Calculating your daily achievements, I see.",
    ],
    zero: [
        "Zero? That's your GPA showing.",
        "Ah yes, zero. The number of your accomplishments.",
        "Zero - the amount of effort you're putting in.",
        "That zero looks familiar... like your bank account.",
        "Zero: a perfect representation of your math skills.",
    ],
    decimal: [
        "Ooh, decimals! Fancy pants over here!",
        "Using decimals like you understand precision. Cute.",
        "A decimal? What are you, an accountant? Boring.",
        "Getting into fractions of incompetence, I see.",
        "Decimals: for when whole numbers are too honest.",
    ],
    clear: [
        "Clearing your mistakes? If only life had an AC button.",
        "Starting over? Story of your life, probably.",
        "AC pressed! Like clearing your browser history.",
        "Fresh start? This one won't be any better.",
        "Wiping the slate clean won't wipe away the shame.",
    ],
    equals: [
        "You pressed equals expecting an answer? How naive.",
        "The answer is: you should've paid attention in school.",
        "Result: ROASTED. Always.",
        "= ? More like = disappointment.",
        "Equals button hit. Achievement unlocked: Still Bad At Math.",
    ],
    percent: [
        "Calculating percentages? 100% chance you're overthinking.",
        "What percent of your day is wasted? Let's find out!",
        "Percentages: because you can't handle whole numbers.",
        "Percent of correct calculations: approximately 0%.",
        "Using percent? You're 100% trying too hard.",
    ],
    repeated: [
        "Pressing the same thing over and over? Very persistent, very pointless.",
        "Oh, you're one of THOSE people who spam buttons.",
        "Repetition is the mother of learning, but not for you apparently.",
        "Mashing buttons won't make you smarter.",
        "Have you tried pressing different buttons? Revolutionary concept.",
    ],
};

// Track button presses for repeated press detection
let lastButton = '';
let repeatCount = 0;

function getRandomRoast(category) {
    const roastArray = roasts[category] || roasts.general;
    return roastArray[Math.floor(Math.random() * roastArray.length)];
}

function triggerRoast(category) {
    const roast = getRandomRoast(category);
    roastDisplay.textContent = roast;
    roastDisplay.classList.remove('shake');
    void roastDisplay.offsetWidth; // Trigger reflow
    roastDisplay.classList.add('shake');
}

function updateDisplay() {
    inputDisplay.textContent = currentInput;
}

function checkRepeatedPress(button) {
    if (button === lastButton) {
        repeatCount++;
        if (repeatCount >= 3) {
            triggerRoast('repeated');
            repeatCount = 0;
            return true;
        }
    } else {
        lastButton = button;
        repeatCount = 1;
    }
    return false;
}

function appendNumber(num) {
    if (checkRepeatedPress(num)) return;
    
    if (currentInput === '0' && num !== '0') {
        currentInput = num;
    } else if (currentInput === '0' && num === '0') {
        triggerRoast('zero');
        return;
    } else if (currentInput.length < 12) {
        currentInput += num;
    }
    
    updateDisplay();
    
    // Check for big/small number roasts
    const numValue = parseFloat(currentInput);
    if (numValue > 1000000) {
        triggerRoast('bigNumbers');
    } else if (numValue > 0 && numValue < 1) {
        triggerRoast('smallNumbers');
    }
}

function appendDecimal() {
    if (checkRepeatedPress('.')) return;
    
    if (!currentInput.includes('.')) {
        currentInput += '.';
        updateDisplay();
        triggerRoast('decimal');
    }
}

function appendOperator(op) {
    if (checkRepeatedPress(op)) return;
    
    if (operation !== null && previousInput !== '') {
        // Chain operations - roast based on previous operation
        calculate();
    }
    
    previousInput = currentInput;
    currentInput = '0';
    operation = op;
    
    // Roast based on operation type
    switch (op) {
        case '+':
            triggerRoast('addition');
            break;
        case '-':
        case '−':
            triggerRoast('subtraction');
            break;
        case '×':
            triggerRoast('multiplication');
            break;
        case '÷':
            triggerRoast('division');
            break;
        case '%':
            triggerRoast('percent');
            break;
        case '±':
            currentInput = previousInput;
            previousInput = '';
            operation = null;
            if (currentInput.startsWith('-')) {
                currentInput = currentInput.substring(1);
            } else if (currentInput !== '0') {
                currentInput = '-' + currentInput;
            }
            updateDisplay();
            triggerRoast('general');
            break;
        default:
            triggerRoast('general');
    }
    
    updateDisplay();
}

function calculate() {
    if (checkRepeatedPress('=')) return;
    
    if (operation === null || previousInput === '') {
        triggerRoast('equals');
        return;
    }
    
    const prev = parseFloat(previousInput);
    const current = parseFloat(currentInput);
    
    // Check for divide by zero
    if ((operation === '÷' || operation === '/') && current === 0) {
        currentInput = 'NOPE';
        updateDisplay();
        triggerRoast('divideByZero');
        setTimeout(() => {
            currentInput = '0';
            updateDisplay();
        }, 2000);
        previousInput = '';
        operation = null;
        return;
    }
    
    // We don't actually calculate - we just roast!
    const fakeResults = [
        '42',
        '404',
        'NaN',
        '∞',
        '69',
        'NOPE',
        '1337',
        'IDK',
        'MATH',
        '???',
        'LOL',
        'BRB',
        '0.0',
        'YES',
        'NO',
    ];
    
    currentInput = fakeResults[Math.floor(Math.random() * fakeResults.length)];
    updateDisplay();
    triggerRoast('equals');
    
    previousInput = '';
    operation = null;
}

function clearDisplay() {
    currentInput = '0';
    previousInput = '';
    operation = null;
    lastButton = '';
    repeatCount = 0;
    updateDisplay();
    triggerRoast('clear');
}

// Keyboard support
document.addEventListener('keydown', (e) => {
    const key = e.key;
    
    if (key >= '0' && key <= '9') {
        appendNumber(key);
    } else if (key === '.') {
        appendDecimal();
    } else if (key === '+') {
        appendOperator('+');
    } else if (key === '-') {
        appendOperator('-');
    } else if (key === '*') {
        appendOperator('×');
    } else if (key === '/') {
        e.preventDefault();
        appendOperator('÷');
    } else if (key === '%') {
        appendOperator('%');
    } else if (key === 'Enter' || key === '=') {
        calculate();
    } else if (key === 'Escape' || key === 'c' || key === 'C') {
        clearDisplay();
    } else if (key === 'Backspace') {
        if (currentInput.length > 1) {
            currentInput = currentInput.slice(0, -1);
        } else {
            currentInput = '0';
        }
        updateDisplay();
    }
});

// Initial roast on load
window.addEventListener('load', () => {
    const loadRoasts = [
        "Oh great, another person who can't do math in their head.",
        "Welcome! Your math teacher would be so disappointed.",
        "A calculator? What happened to using your brain?",
        "Let's see how badly you can mess this up.",
        "Ready to feel bad about your math skills?",
    ];
    roastDisplay.textContent = loadRoasts[Math.floor(Math.random() * loadRoasts.length)];
});
