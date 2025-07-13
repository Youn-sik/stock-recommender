FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    wget \
    curl \
    unzip \
    # Chrome dependencies
    chromium \
    chromium-driver \
    && rm -rf /var/lib/apt/lists/*

# Install Korean language support for konlpy
RUN apt-get update && apt-get install -y \
    default-jdk \
    && rm -rf /var/lib/apt/lists/*

# Set Java environment
ENV JAVA_HOME=/usr/lib/jvm/default-java

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Download NLTK data
RUN python -c "import nltk; nltk.download('punkt'); nltk.download('vader_lexicon')"

# Copy application code
COPY . .

# Create logs directory
RUN mkdir -p logs

# Set Chrome options for headless mode
ENV CHROME_OPTIONS="--headless --no-sandbox --disable-dev-shm-usage"

CMD ["python", "main.py"]