# Quizine

A flexible and extensible quiz engine written in Go. Quizine is designed to provide a core set of functionality and data handling, allowing developers to easily integrate quizzing functionality into their applications.

## Core Concepts

* **Quiz:** Defines the structure and rules of a specific test or assessment, containing a collection of questions and configuration settings.
* **User:** Represents an individual or entity taking a quiz. This structure holds identifying information and tracks their progress or history.
* **Quiz Attempt:** Records a single instance of a `User` (or group) taking a specific `Quiz`. It tracks the state during the attempt, including the questions presented and the user's responses.
* **Question Attempt:** Records the specific details and the user's response for one individual question presented during a `QuizAttempt`.
* **Question Result:** Stores the evaluated outcome for a single `QuestionAttempt`, including correctness, score awarded, and feedback.
* **Quiz Result:** Stores the final outcome and details of a completed `QuizAttempt`, including the score, grading breakdown, and potentially feedback.
* **Question:** Represents a single item in a quiz designed to assess a user's knowledge or understanding. It comprises a prompt, an input mechanism for the answer, and the data required to determine correctness (as detailed in the Question Types section).
* **Prompt:** The text, media, or stimulus presented to the user that poses the question or task.
* **Input Mechanism (or Answer Format):** The method by which the user provides their answer, determined by the question type.
* **Answer Data:** The information associated with the question that allows for its resolution or evaluation (e.g., correct responses, options, sequences).

### Example Concept Flow
1. A `Quiz` exists called 'History 101' that contains 5 `Question`s, each of which contains a `Prompt`, `Input Mechanism`, and `Answer Data`.
2. A `User` starts the `Quiz`.
3. A new `Quiz Attempt` record is created, linking the `User` to the `Quiz` and recording the `StartTime`. The specific `Question`s for this attempt (potentially shuffled) are determined and stored.
4. For each `Question` in the `Quiz Attempt`:   
4.1 The `Prompt` and `Input Mechanism` of the current `Question` are presented to the `User`.  
4.2 The `User` provides input via the `Input Mechanism`.  
4.3 A `Question Attempt` record is created, linked to the `Quiz Attempt` and the specific `Question`, storing the `User Answer` and `TimeSpent` on that question.  
4.4 The `User Answer` from the `Question Attempt` is evaluated against the `Answer Data` of the original `Question`.  
4.5 A `Question Result` record is created, linked to the `Question Attempt`, storing the evaluation outcome (`IsCorrect`, `ScoreAwarded`, `Feedback`).  
5. Once the `User` has answered all questions or the `TimeLimit` of the `Quiz Attempt` is reached, the `EndTime` is recorded, and the `Status` is set to "Completed" or "Expired".
6. A final `Quiz Result` record is generated, linked to the completed `Quiz Attempt`. This involves aggregating the `ScoreAwarded` from all `Question Result`s to calculate the total `Score` and `Percentage`, determining if the user `Passed`, and storing the `CompletionTime`.

---

# Data

Quizine is built with a flexible data layer. By default, it uses **SQLite** for persistence, offering a simple, file-based database solution out-of-the-box. However, the data interfaces are designed to be easily configurable, allowing you to implement support for other databases (like PostgreSQL, MySQL, etc.) or even different storage mechanisms as needed.

## Data Structures

Here's a breakdown of the key data structures within Quizine. These structures define the information needed to manage quizzes, users, and their interactions with the quizzes.

### Quiz

The `Quiz` structure serves as the blueprint for a particular test. It defines the content and rules that govern a user's attempt.

* **Description:** A container for a set of `Question`s, along with metadata and configuration that determine how the quiz is presented, taken, and scored.
* **Data Fields:**
    * `ID`: A globally unique identifier for the quiz.
    * `Title`: The name of the quiz (e.g., "Introduction to Go", "History Quiz").
    * `Description`: A brief explanation of the quiz's topic or purpose.
    * `Instructions`: Instructions for the user taking the quiz.
    * `Questions`: An ordered collection of the `Question`s included in this quiz, referenced by their IDs.
    * `TimeLimit`: The maximum duration allowed for completing the quiz (e.g., in minutes or seconds). Can be optional (0 for no limit).
    * `PassingScore`: The minimum score or percentage required to pass the quiz. Optional too (0 for a pass regardless of score).
    * `ScoringMethod`: Defines how the quiz is scored (e.g., points per question, penalty for incorrect answers).
    * `ShuffleQuestions`: Boolean for the order of questions should be randomized for each attempt.
    * `ShuffleAnswers`: Boolean for the order of answer options (for types like Multi-Choice) should be randomized.
    * `CreatedAt`, `UpdatedAt`: Timestamps for tracking when the quiz was created and last modified.
    * `Status`: The current state of the quiz (e.g., "Draft", "Published", "Archived").

---

### User

The `User` represents the individual or entity interacting with the Quizine engine, typically a person taking a quiz.

* **Description:** Holds identifying information for a participant and serves as a link to track their quiz attempts and results.
* **Data Fields:**
    * `ID`: A globally unique identifier for the user.
    * `Email`: A unique identifier for login/identification purposes and the email for contact.
    * `DisplayName`: The name displayed to the user or in results.
    * `CreatedAt`, `UpdatedAt`: Timestamps.
    * `Roles`: Represents the level of access within the engine to allow or restrict certain permissions.

---

### QuizAttempt

The `QuizAttempt` structure records the details of a single instance of a user taking a specific quiz from start to finish. It's the live record of the user's interaction.

* **Description:** Captures the state and progress of a user's journey through a particular quiz, including the specific questions presented and the user's responses before final scoring.
* **Data Fields:**
    * `ID`: A globally unique identifier for this specific attempt.
    * `QuizID`: A reference to the `Quiz` that was attempted.
    * `UserIDs`: A reference to the `User`(s) who made the attempt.
    * `StartTime`: Timestamp when the user started the attempt.
    * `EndTime`: Timestamp when the user finished or the time limit expired.
    * `Status`: The current state of the attempt (e.g., "In Progress", "Completed", "Expired").
    * `UserAnswers`: A collection storing the user's response(s) for each question in this attempt, stored in a QuestionAttempt.
    * `QuestionsPresented`: An ordered list of references to Question IDs indicating the specific questions shown to the user in this attempt, especially relevant if questions are shuffled or drawn from a pool.
    * *Timer State (if applicable):* Data to track the remaining time during the attempt.

---

### Question Attempt

The `QuestionAttempt` structure stores the specific details and the user's response for one individual question presented during a `QuizAttempt`.

* **Description:** Records the user's interaction with a single `Question` within the scope of a `QuizAttempt`. This structure holds the user's submitted answer and not the outcome.
* **Data Fields:**
    * `ID`: A unique identifier for this specific question attempt.
    * `AttemptID`: A reference to the `QuizAttempt` this question attempt belongs to.
    * `QuestionID`: A reference to the `Question` that was presented.
    * `UserAnswer`: The data submitted by the user as their answer to this specific question. The format and type of this data will vary significantly based on the `Question`'s `Input Mechanism` (e.g., an index for Multi-Choice, a string for Fill-in, coordinates for Hotspot).
    * `TimeSpent`: The duration the user spent on this specific question during the attempt.

---

### Question Result

The `QuestionResult` structure stores the evaluated outcome for a single `QuestionAttempt`. It represents the grading and feedback for that specific question response.

* **Description:** Holds the result of evaluating a user's `UserAnswer` from a `QuestionAttempt` against the correct criteria defined in the `Question`. This structure contains the outcome of the grading process for a single question.
* **Data Fields:**
    * `ID`: A unique identifier for this specific question result.
    * `QuestionAttemptID`: A reference to the `QuestionAttempt` this result is based on.
    * `IsCorrect`: A boolean indicating if the user's answer was evaluated as correct.
    * `ScoreAwarded`: The points or partial points awarded for this specific question based on the evaluation and the quiz's scoring rules.
    * `Feedback`: Any specific feedback generated for the user's answer on this question (e.g., explaining why it was incorrect, providing the correct answer).
    * `EvaluatedAt`: Timestamp when the question was evaluated (this might happen instantly or upon completion of the attempt).

---

### Quiz Result

The `QuizResult` structure provides a summary and breakdown of a completed `QuizAttempt`. It contains the evaluated outcome.

* **Description:** Stores the final score and details of how a user performed on a completed `QuizAttempt`, providing a persistent record of their achievement.
* **Data Fields:**
    * `ID`: A unique identifier for the result.
    * `AttemptID`: A reference to the `QuizAttempt` this result is based on.
    * `QuizID`: A reference to the `Quiz` that was attempted.
    * `UserID`: A reference to the `User` who made the attempt.
    * `Score`: The final numerical score obtained by the user.
    * `Percentage`: The score expressed as a percentage.
    * `Passed`: Boolean indicating whether the user achieved the `PassingScore` for the quiz.
    * `CompletionTime`: The total time taken to complete the quiz (derived from `EndTime` - `StartTime` in `QuizAttempt`).
    * `GradingDetails`: A detailed breakdown of the grading, potentially including points awarded/deducted for each question, or indicating which specific answers were correct/incorrect.
    * `CompletedAt`: Timestamp when the result was generated (usually matches the `EndTime` of the attempt).
    * *Feedback (Optional):* General feedback or per-question feedback.

---

### Question
TODO

### Question Types

Quizine supports a variety of question types to cater to different assessment needs. Each type has specific data requirements:

#### Multi-Choice

* **Description:** Presents a question prompt and a list of possible answers, from which the user must select **only one** correct option.
* **Data Structure Needs:**
    * `Prompt`: The question text.
    * `Options`: A list of possible answer strings.
    * `CorrectAnswerIndex`: The index (position) of the single correct answer within the `Options` list.

#### True/False

* **Description:** Presents a statement, and the user must determine if the statement is true or false.
* **Data Structure Needs:**
    * `Prompt`: The statement to be evaluated.
    * `CorrectAnswer`: A boolean value (`true` or `false`) indicating the correct evaluation of the statement.

#### Multi-Response

* **Description:** Similar to Multi-Choice, but the user must select **all** correct options from the provided list. There can be one, multiple, or even zero correct answers depending on the question design.
* **Data Structure Needs:**
    * `Prompt`: The question text.
    * `Options`: A list of possible answer strings.
    * `CorrectAnswerIndices`: A list of indices corresponding to all the correct answers within the `Options` list.

#### Matching

* **Description:** Presents two lists of items. The user must create pairs by matching an item from the first list with its corresponding item from the second list.
* **Data Structure Needs:**
    * `Prompt`: Instructions for the matching task.
    * `List1Items`: The items in the first list.
    * `List2Items`: The items in the second list.
    * `CorrectPairings`: A list of pairs of indices representing the correct matches between `List1Items` and `List2Items`.

#### Ordering/Sequence

* **Description:** Provides a set of items that are initially presented in an unordered or random manner. The user must arrange these items into a specific correct sequence.
* **Data Structure Needs:**
    * `Prompt`: Instructions for the ordering task.
    * `Items`: The list of items to be ordered.
    * `CorrectOrderIndices`: A list of indices representing the correct sequence of the `Items` list.

#### Drag-and-Drop

* **Description:** Requires the user to interact by moving a digital object (text, image) from one location to a designated target area on the screen or within text.
* **Variations:**
    * **Drag-and-Drop into Text (Cloze):** Drag words/phrases into blanks in a sentence.
    * **Drag-and-Drop onto Image:** Drag markers/labels onto specific points on an image.
* **Data Structure Needs (Example for Drag-and-Drop into Text):**
    * `Prompt`: The text containing the blanks.
    * `DroppableItems`: The list of words or phrases to be dragged.
    * `CorrectDropTargets`: Definitions of where each `DroppableItem` should be placed within the `Prompt`.

#### Select Missing (Dropdown)

* **Description:** Presents text with one or more blanks. Each blank has a dropdown menu from which the user must select the correct word or phrase to complete the text. This is similar to Fill-in but provides options.
* **Data Structure Needs:**
    * `PromptTemplate`: The text with placeholders for the dropdowns.
    * `Dropdowns`: A list where each item defines the options for a specific dropdown and indicates the correct option within that list.

#### Fill-in

* **Description:** Requires the user to type a word, phrase, or sentence directly into a text input field to complete a statement or answer a question.
* **Data Structure Needs:**
    * `Prompt`: The statement or question.
    * `CorrectAnswers`: A list of acceptable text strings. (Allows for variations like different capitalization or minor synonyms).
    * `CaseSensitive`: Boolean indicating if the answer matching should be case-sensitive.

#### Numeric

* **Description:** Requires the user to input a numerical value as the answer. This is ideal for mathematical or scientific questions with a quantitative answer.
* **Data Structure Needs:**
    * `Prompt`: The question requiring a numerical answer.
    * `CorrectAnswer`: The exact numerical value.
    * `Tolerance` (Optional): An acceptable range around the `CorrectAnswer` to allow for slight variations (e.g., floating-point inaccuracies).
    * `AllowRange` (Optional): Define a minimum and maximum acceptable value instead of a single correct answer with tolerance.

#### Hotspot

* **Description:** Displays an image, and the user must click on a specific predefined area (or areas) within the image to provide the answer.
* **Data Structure Needs:**
    * `Prompt`: The question asking the user to identify something in the image.
    * `Image`: Reference to the image file or data.
    * `CorrectAreas`: Definitions of the clickable regions in the image that are considered correct (e.g., coordinates or shapes).

#### Audio

* **Description:** Requires the user to provide an answer by recording and submitting an audio file. This is useful for language proficiency tests or capturing verbal responses.
* **Data Structure Needs:**
    * `Prompt`: The question or instruction for the audio response.
    * *Evaluation:* Typically requires manual grading or integration with speech-to-text and potentially AI analysis. The data structure might primarily store the prompt and a placeholder/reference for the submitted audio.

#### Video

* **Description:** Requires the user to provide an answer by recording and submitting a video file. Suitable for demonstrating physical actions, presentations, or capturing more complex responses.
* **Data Structure Needs:**
    * `Prompt`: The question or instruction for the video response.
    * *Evaluation:* Typically requires manual grading. The data structure might store the prompt and a placeholder/reference for the submitted video.

---



