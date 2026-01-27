# 22 - AI Capabilities

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 3

---

## 1. Overview

The AI Capabilities module provides intelligent insights, predictions, and automation using on-premise AI models for data privacy. All AI features are optional and modular.

---

## 2. AI Architecture

### 2.1 On-Premise AI Stack

```
┌─────────────────────────────────────────────────────────────┐
│                    APPLICATION LAYER                         │
│              (School ERP Backend - Go)                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    AI SERVICE LAYER                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Prediction  │  │   NLP       │  │ Recommend   │         │
│  │   Engine    │  │  Engine     │  │   Engine    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    MODEL LAYER                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Local LLM (Llama 3, Mistral, Phi)                  │    │
│  │  Runs on school server - No data leaves premises    │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  ML Models (scikit-learn, XGBoost)                  │    │
│  │  Trained on school's historical data               │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Privacy-First Approach

| Principle | Implementation |
|-----------|----------------|
| Data Stays Local | All processing on school server |
| No Cloud AI | No OpenAI, Google AI dependencies |
| Explainable AI | Show reasoning for predictions |
| Opt-in Features | Schools enable AI features explicitly |
| Audit Trail | Log all AI decisions |

---

## 3. Academic AI Features

### 3.1 Student Performance Prediction

```
Student Risk Assessment | Aarav Sharma | Class 10-A

PREDICTION: Moderate Risk of Declining Performance

Risk Score: 65/100 (Moderate)

CONTRIBUTING FACTORS:
├─ Attendance dropped from 95% to 88% (last 30 days)
├─ Math quiz scores declining trend (-15% over 3 quizzes)
├─ Homework completion rate dropped to 70%
└─ Less engagement in digital classroom (watch time down)

RECOMMENDED ACTIONS:
1. Schedule parent-teacher meeting
2. Provide extra math practice materials
3. Assign peer study partner
4. Monitor for next 2 weeks

SIMILAR PATTERNS:
5 students showed similar patterns last year:
- 3 improved with intervention
- 2 declined further

[Send Alert to Class Teacher] [Schedule PTM]
```

### 3.2 Dropout Risk Prediction

```
Dropout Risk Analysis | January 2026

HIGH RISK STUDENTS (Score > 80):
| Student        | Class | Risk Score | Key Factors            |
|----------------|-------|------------|------------------------|
| Rahul Kumar    | 9-B   | 92         | Attendance, Fee Dues   |
| Priya Singh    | 8-A   | 85         | Grades, Engagement     |
| Amit Verma     | 10-A  | 82         | Attendance, Behavior   |

EARLY WARNING (Score 60-80): 12 students
MONITOR (Score 40-60): 25 students

Model Accuracy: 87% (based on last year's predictions)
```

### 3.3 Personalized Learning Recommendations

```
Learning Recommendations | Aarav Sharma

BASED ON YOUR PERFORMANCE:

Mathematics:
├─ Strength: Algebra (85% avg)
├─ Weakness: Geometry (62% avg)
└─ Recommended:
   • Watch: "Circle Theorems Explained" (12 min)
   • Practice: Geometry Quiz Set 3
   • Review: Chapter 8 Notes

Science:
├─ Strength: Physics (78% avg)
├─ Weakness: Chemistry (65% avg)
└─ Recommended:
   • Watch: "Chemical Equations Basics" (15 min)
   • Practice: Balancing Equations Worksheet

STUDY PLAN FOR THIS WEEK:
| Day | Focus Area        | Time | Resources         |
|-----|-------------------|------|-------------------|
| Mon | Geometry          | 1 hr | Video + Practice  |
| Tue | Chemistry         | 1 hr | Notes + Quiz      |
| Wed | Mixed Practice    | 1 hr | Past Papers       |
```

---

## 4. Administrative AI Features

### 4.1 Fee Default Prediction

```
Fee Default Prediction | February 2026

PREDICTED DEFAULTS (High Probability > 70%):

| Student        | Class | Due Amount | Probability | Factors         |
|----------------|-------|------------|-------------|-----------------|
| Rahul Verma    | 7-A   | ₹25,000    | 85%         | History, Income |
| Meera Patel    | 9-B   | ₹18,000    | 78%         | Pattern, Delay  |
| Arjun Singh    | 10-A  | ₹32,000    | 72%         | Previous Default|

RECOMMENDED ACTIONS:
├─ Send personalized reminder 7 days before due date
├─ Offer installment plan to high-risk families
└─ Schedule finance counseling

TOTAL AT RISK: ₹3.2L from 15 students

[Generate Collection Strategy] [Send Reminders]
```

### 4.2 Admission Demand Forecasting

```
Admission Forecast | 2026-27

PREDICTED APPLICATIONS BY CLASS:
| Class | Last Year | Predicted | Confidence |
|-------|-----------|-----------|------------|
| LKG   | 150       | 165       | High       |
| UKG   | 120       | 125       | High       |
| 1     | 80        | 95        | Medium     |
| 2     | 45        | 50        | Medium     |
| 6     | 60        | 70        | High       |

PEAK APPLICATION PERIODS:
├─ January (40% of applications)
├─ February (30% of applications)
└─ March (20% of applications)

RECOMMENDATIONS:
├─ Increase LKG seats from 60 to 70
├─ Start marketing campaign in December
└─ Prepare for 15% more applications than last year
```

### 4.3 Staff Workload Optimization

```
Workload Analysis | Teaching Staff

OVERLOADED (>30 periods/week):
| Teacher        | Current | Optimal | Subjects        |
|----------------|---------|---------|-----------------|
| Mr. Kumar      | 35      | 28      | Math, Physics   |
| Ms. Sharma     | 32      | 28      | English         |

UNDERUTILIZED (<20 periods/week):
| Teacher        | Current | Capacity | Subjects       |
|----------------|---------|----------|----------------|
| Mr. Patel      | 18      | 28       | Science        |
| Ms. Gupta      | 16      | 28       | Hindi          |

OPTIMIZATION SUGGESTIONS:
├─ Transfer 3 Math periods from Mr. Kumar to available staff
├─ Assign additional Science sections to Mr. Patel
└─ Potential savings: 1 temporary staff position

[Generate Optimized Timetable]
```

---

## 5. AI-Assisted Content

### 5.1 Auto-Generated Summaries

```
Class Recording Summary | Quadratic Equations | 22-Jan-2026

TOPICS COVERED:
1. Introduction to Quadratic Equations (0:00 - 5:30)
2. Standard Form ax² + bx + c = 0 (5:30 - 12:00)
3. Solving by Factorization (12:00 - 25:00)
4. Completing the Square Method (25:00 - 38:00)
5. Practice Problems (38:00 - 45:00)

KEY CONCEPTS:
├─ A quadratic equation has degree 2
├─ Standard form requires a ≠ 0
├─ Two methods taught: Factorization and Completing Square
└─ Discriminant determines nature of roots

FORMULAS MENTIONED:
├─ x = (-b ± √(b²-4ac)) / 2a
└─ Completing square: (x + b/2a)² = (b²-4ac)/4a²

AUTO-GENERATED QUIZ: [Create Quiz from Content]
```

### 5.2 Question Generation

```
AI Question Generator | Chapter 5: Quadratic Equations

GENERATED QUESTIONS:

MCQ (Easy):
Q1. Which of the following is a quadratic equation?
    a) x + 2 = 0
    b) x² + 3x + 2 = 0 ✓
    c) x³ + x = 0
    d) 2x = 4

MCQ (Medium):
Q2. The roots of x² - 5x + 6 = 0 are:
    a) 2, 3 ✓
    b) -2, -3
    c) 1, 6
    d) -1, -6

Short Answer (Hard):
Q3. Solve using completing the square: x² + 6x + 5 = 0
    Expected Answer: x = -1 or x = -5

[Generate More] [Add to Question Bank] [Create Quiz]
```

### 5.3 Multilingual Content Support

```
Content Translation | Science Notes

ORIGINAL (English):
"Photosynthesis is the process by which plants convert
sunlight into chemical energy (glucose)."

TRANSLATED (Hindi):
"प्रकाश संश्लेषण वह प्रक्रिया है जिसके द्वारा पौधे
सूर्य के प्रकाश को रासायनिक ऊर्जा (ग्लूकोज) में बदलते हैं।"

TRANSLATED (Marathi):
"प्रकाशसंश्लेषण ही प्रक्रिया आहे ज्याद्वारे वनस्पती
सूर्यप्रकाशाचे रासायनिक ऊर्जेत (ग्लुकोज) रूपांतर करतात."

[Translate Full Chapter] [Review & Edit]
```

---

## 6. Anomaly Detection

### 6.1 Attendance Anomalies

```
Attendance Anomaly Alert

DETECTED: Unusual absence pattern in Class 8-B

Pattern: 8 students absent on same 3 days (Mon, Wed, Fri)
Normal rate: 3-4 absences/day
Detected rate: 8 absences on specific days

POSSIBLE CAUSES:
├─ Tuition class conflict
├─ Transport issue on specific route
└─ Group activity outside school

RECOMMENDED ACTION:
├─ Investigate common factors
├─ Contact parents of affected students
└─ Check if pattern continues

[Investigate] [Dismiss]
```

### 6.2 Grade Anomalies

```
Grading Anomaly Alert

DETECTED: Unusual grade distribution in Math Unit Test

Class: 10-A | Teacher: Mr. Kumar
Expected distribution: Normal curve
Actual: Bimodal (high at 80-90 and 40-50)

ANALYSIS:
├─ 60% students scored 80-90 (unusually high)
├─ 30% students scored 40-50 (unusually low)
├─ Possible test leak or teaching gap

RECOMMENDATION:
├─ Review test security
├─ Analyze question-wise performance
├─ Consider remedial class for low scorers

[Review Details] [Dismiss]
```

---

## 7. AI Configuration

### 7.1 Feature Toggles

```
AI Features Configuration

| Feature                    | Status   | Data Required |
|----------------------------|----------|---------------|
| Performance Prediction     | Enabled  | 1 year        |
| Dropout Risk              | Enabled  | 2 years       |
| Fee Default Prediction    | Disabled | 1 year        |
| Learning Recommendations  | Enabled  | 3 months      |
| Content Summarization     | Enabled  | -             |
| Question Generation       | Enabled  | -             |
| Anomaly Detection         | Enabled  | 6 months      |

Model Last Updated: 15-Jan-2026
Next Scheduled Training: 01-Feb-2026
```

### 7.2 Model Training

- Models retrained monthly with new data
- Historical data anonymized for training
- School can opt-out of specific features
- Explainability reports available

---

## 8. API Endpoints

```
# Predictions
GET    /api/v1/ai/students/{id}/risk        # Student risk
GET    /api/v1/ai/students/dropout-risk     # Dropout predictions
GET    /api/v1/ai/fees/default-prediction   # Fee defaults
GET    /api/v1/ai/admissions/forecast       # Admission forecast

# Recommendations
GET    /api/v1/ai/students/{id}/learning    # Learning plan
GET    /api/v1/ai/staff/workload            # Workload analysis

# Content
POST   /api/v1/ai/content/summarize         # Summarize content
POST   /api/v1/ai/content/questions         # Generate questions
POST   /api/v1/ai/content/translate         # Translate content

# Anomalies
GET    /api/v1/ai/anomalies                 # Active anomalies
POST   /api/v1/ai/anomalies/{id}/dismiss    # Dismiss anomaly

# Configuration
GET    /api/v1/ai/config                    # AI settings
PUT    /api/v1/ai/config                    # Update settings
```

---

## 9. Hardware Requirements

### 9.1 On-Premise AI Server

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 8 cores | 16+ cores |
| RAM | 32 GB | 64 GB |
| GPU | Optional | NVIDIA RTX 3080+ |
| Storage | 500 GB SSD | 1 TB NVMe |

### 9.2 Without GPU (CPU-only)

- Use quantized models (GGUF format)
- Slower inference but functional
- Suitable for smaller schools

---

## 10. Related Documents

- [21-analytics-dashboards.md](./21-analytics-dashboards.md) - Analytics
- All module documents provide data for AI
- [index.md](./index.md) - Main PRD index

---

**Previous**: [21-analytics-dashboards.md](./21-analytics-dashboards.md)
**End of PRD Documents**
