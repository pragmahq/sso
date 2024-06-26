# Nexus

## Overview

Nexus, inspired by the central hub in [Demon's Souls](https://en.wikipedia.org/wiki/Demon%27s_Souls), serves as the centralized platform for internal Pragma users. This Next.js application provides essential functionality for our team members, streamlining their daily workflows and enhancing productivity.

## Key Features

1. **Global Profile Management**: Users can easily update and maintain their comprehensive profile information, ensuring consistency across all integrated systems.

2. **Unified Authentication**: Nexus provides a single sign-on (SSO) solution, allowing users to access all connected Pragma applications with one set of credentials.

## Development Instructions

To set up and run the Nexus project locally, follow these steps:

1. **Prerequisites**:
   - Ensure you have Node.js (version 14 or later) installed on your machine.
   - Familiarity with React and Next.js is recommended.

2. **Clone the Repository**:
   ```
   git clone https://github.com/pragma/nexus.git
   cd nexus
   ```

3. **Install Dependencies**:
   ```
   npm install
   ```

4. **Environment Setup**:
   - Copy the `.env.example` file to `.env.local`.
   - Update the environment variables with the necessary API keys and configuration details.

5. **Run the Development Server**:
   ```
   npm run dev
   ```

6. **Access the Application**:
   Open your browser and navigate to `http://localhost:3000` to view the Nexus app.

7. **Building for Production**:
   When you're ready to deploy, use the following commands:
   ```
   npm run build
   npm start
   ```

Remember to adhere to our coding standards and contribute to the documentation as you develop new features or make changes to existing ones.
