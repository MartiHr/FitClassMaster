-- PostgreSQL Data Seed Script
-- 1. Clean up old data (Optional, safer for fresh DBs)
TRUNCATE TABLE "session_logs" CASCADE;
TRUNCATE TABLE "workout_exercises" CASCADE;
TRUNCATE TABLE "workout_sessions" CASCADE;
TRUNCATE TABLE "messages" CASCADE;
TRUNCATE TABLE "conversations" CASCADE;
TRUNCATE TABLE "enrollments" CASCADE;
TRUNCATE TABLE "classes" CASCADE;
TRUNCATE TABLE "workout_plans" CASCADE;
TRUNCATE TABLE "exercises" CASCADE;
TRUNCATE TABLE "users" CASCADE;

-- 2. Insert Data
-- Users
INSERT INTO "users" ("id", "first_name", "last_name", "email", "password", "role", "created_at", "updated_at") VALUES 
(1, 'Gosho', 'Goshev', 'george@abv.bg', '$2a$10$.LPU3FStRPIPSQqsp8IaLeZMIr175E8HQfwlEpSBM3n3qUiyl6oFC', 'member', '2026-01-18 14:08:21+02', '2026-02-03 16:31:48+02'),
(2, 'Pesho', 'peshev', 'pesho@abv.bg', '$2a$10$zORC2t3koKo1DOvtOffW9OKpSlsUZvmq8IyhVqjWZFqyMBeK78TH2', 'member', '2026-01-18 15:51:10+02', '2026-02-03 19:55:41+02'),
(3, 'Test', 'Testov', 'trainer@test.com', '$2a$10$BWQ1fDR09/bOFqZX3kzNSep9YeMzGAWLaFTl3KcXYmuWeTCI62U8.', 'trainer', '2026-01-22 15:44:40+02', '2026-01-22 15:44:40+02'),
(4, 'Admin', 'Adminov', 'admin@admin.com', '$2a$10$bkBBURpp.AvjZJe33nKs.OTC82/MPwZWvsgY0G.tKmIRHACqKbVrm', 'admin', '2026-02-03 16:28:42+02', '2026-02-03 16:28:42+02');

-- Exercises
INSERT INTO "exercises" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "muscle_group", "equipment", "video_url") VALUES 
(1, '2026-01-22 16:05:34+02', '2026-01-22 16:05:34+02', NULL, 'Bench Press', 'Pull up and down.', 'Chest', 'Barbell Machine', 'https://www.youtube.com/shorts/Z3kjNh5UKqI'),
(2, '2026-01-22 16:05:34+02', '2026-01-22 16:05:34+02', '2026-01-22 16:08:58+02', 'Bench Press', 'Pull up and down.', 'Chest', 'Barbell Machine', 'https://www.youtube.com/shorts/Z3kjNh5UKqI'),
(3, '2026-01-22 16:17:11+00', '2026-01-22 16:17:11+00', NULL, 'Barbell Squat', 'Stand with feet shoulder-width apart. Lower your hips back and down as if sitting in a chair, keeping your chest up.', 'Legs', 'Barbell', 'https://www.youtube.com/watch?v=ultWZbGWL54'),
(4, '2026-01-22 16:17:11+00', '2026-01-22 16:17:11+00', NULL, 'Push-Up', 'Start in a plank position. Lower your body until your chest nearly touches the floor, then push back up.', 'Chest', 'Bodyweight', 'https://www.youtube.com/watch?v=IODxDxX7oi4'),
(5, '2026-01-22 16:17:11+00', '2026-01-22 16:17:11+00', NULL, 'Pull-Up', 'Hang from a bar with palms facing away. Pull your body up until your chin is over the bar.', 'Back', 'Pull-up Bar', 'https://www.youtube.com/watch?v=eGo4IYlbE5g'),
(6, '2026-01-22 16:17:11+00', '2026-01-22 16:17:11+00', NULL, 'Dumbbell Shoulder Press', 'Sit or stand with a dumbbell in each hand at shoulder height. Press weights overhead until arms are extended.', 'Shoulders', 'Dumbbells', 'https://www.youtube.com/watch?v=qEwK6KpUeCQ'),
(7, '2026-01-22 16:17:11+00', '2026-01-22 16:17:11+00', NULL, 'Plank', 'Hold a push-up position but resting on your forearms. Keep your body in a straight line from head to heels.', 'Core', 'Bodyweight', 'https://www.youtube.com/watch?v=pSHjTRCQxIw'),
(8, '2026-01-22 16:17:11+00', '2026-01-22 16:17:11+00', NULL, 'Deadlift', 'Stand with feet hip-width apart. Hinge at hips to grab the bar, keep back flat, and lift by extending hips and knees.', 'Legs', 'Barbell', 'https://www.youtube.com/watch?v=op9kVnSso6Q'),
(9, '2026-01-22 16:17:11+00', '2026-01-22 16:17:11+00', NULL, 'Burpees', 'Drop into a squat, kick feet back to plank, do a push-up, jump feet back in, and jump up explosively.', 'Cardio', 'Bodyweight', 'https://www.youtube.com/watch?v=auBLPEO8Twy'),
(10, '2026-02-03 18:45:42+02', '2026-02-03 18:45:42+02', '2026-02-04 11:47:50+02', 'DELETE', '', 'Chest', '', '');

-- Workout Plans
INSERT INTO "workout_plans" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "trainer_id") VALUES 
(1, '2026-01-22 17:07:26+02', '2026-01-22 17:07:26+02', NULL, 'High intensity leg workout', 'Burn fat', 3),
(2, '2026-01-22 17:11:42+00', '2026-02-03 13:53:24+02', NULL, 'Full Body Beginner', 'A balanced routine hitting all major muscle groups. Great for starting out.', 3),
(3, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 'HIIT Cardio Blast', 'High intensity interval training to burn fat and improve endurance.', 3),
(4, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 'Upper Body Strength', 'Focus on building mass in shoulders, chest, and back.', 3),
(5, '2026-01-22 17:34:47+02', '2026-01-22 17:34:47+02', NULL, 'A', 'A', 3),
(6, '2026-02-03 18:49:09+02', '2026-02-03 18:49:09+02', '2026-02-03 18:58:32+02', 'DELETE', '', 3);

-- Workout Exercises (Linking tables)
INSERT INTO "workout_exercises" ("id", "created_at", "updated_at", "deleted_at", "workout_plan_id", "exercise_id", "sets", "reps", "duration", "notes", "order") VALUES 
(1, '2026-01-22 17:07:26+02', '2026-01-22 17:07:26+02', NULL, 1, 3, 3, 10, 0, '', 1),
(6, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 3, 7, 4, 15, 0, 'Explosive movement', 1),
(7, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 3, 1, 4, 20, 0, 'Bodyweight squats, fast pace', 2),
(8, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 3, 2, 4, 15, 0, 'Minimal rest', 3),
(9, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 4, 4, 4, 8, 0, 'Heavy weight', 1),
(10, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 4, 3, 4, 8, 0, 'Weighted if possible', 2),
(11, '2026-01-22 17:11:42+00', '2026-01-22 17:11:42+00', NULL, 4, 6, 3, 5, 0, 'Focus on form, heavy', 3),
(12, '2026-01-22 17:34:47+02', '2026-01-22 17:34:47+02', NULL, 5, 9, 3, 10, 0, '', 1),
(13, '2026-02-03 13:53:24+02', '2026-02-03 13:53:24+02', NULL, 2, 1, 3, 10, 0, 'Focus on depth', 1),
(14, '2026-02-03 13:53:24+02', '2026-02-03 13:53:24+02', NULL, 2, 1, 3, 12, 0, 'Keep core tight', 2),
(15, '2026-02-03 13:53:24+02', '2026-02-03 13:53:24+02', NULL, 2, 3, 3, 8, 0, 'Use assistance if needed', 3),
(16, '2026-02-03 13:53:24+02', '2026-02-03 13:53:24+02', NULL, 2, 5, 3, 10, 0, 'Hold for 45 seconds', 4),
(17, '2026-02-03 13:53:24+02', '2026-02-03 13:53:24+02', NULL, 2, 9, 1, 1, 0, '', 5),
(18, '2026-02-03 18:49:09+02', '2026-02-03 18:49:09+02', NULL, 6, 10, 3, 10, 0, '', 1);

-- Workout Sessions
INSERT INTO "workout_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "workout_plan_id", "start_time", "end_time", "status") VALUES 
(1, '2026-01-22 19:51:11+02', '2026-01-22 19:51:11+02', NULL, 3, 3, '2026-01-22 19:51:11+02', NULL, 'completed'),
(2, '2026-01-22 19:53:07+02', '2026-01-22 19:53:07+02', NULL, 3, 1, '2026-01-22 19:53:07+02', NULL, 'completed'),
(3, '2026-01-22 20:08:51+02', '2026-01-22 20:08:51+02', NULL, 3, 1, '2026-01-22 20:08:51+02', NULL, 'completed'),
(4, '2026-01-22 20:08:59+02', '2026-01-22 20:08:59+02', NULL, 3, 2, '2026-01-22 20:08:59+02', NULL, 'completed'),
(5, '2026-01-22 20:09:22+02', '2026-01-22 20:09:22+02', NULL, 1, 2, '2026-01-22 20:09:22+02', NULL, 'completed'),
(6, '2026-01-22 20:09:22+02', '2026-01-22 20:09:22+02', NULL, 1, 2, '2026-01-22 20:09:22+02', NULL, 'completed'),
(7, '2026-01-22 20:10:50+02', '2026-01-22 20:10:50+02', NULL, 1, 2, '2026-01-22 20:10:50+02', NULL, 'completed'),
(8, '2026-01-22 20:10:50+02', '2026-01-22 20:10:50+02', NULL, 1, 2, '2026-01-22 20:10:50+02', NULL, 'completed'),
(9, '2026-01-22 20:13:32+02', '2026-01-22 20:13:32+02', NULL, 1, 2, '2026-01-22 20:13:32+02', NULL, 'completed'),
(10, '2026-01-22 22:15:57+02', '2026-01-22 22:15:57+02', NULL, 1, 2, '2026-01-22 22:15:57+02', NULL, 'completed'),
(11, '2026-01-22 22:17:12+02', '2026-01-22 22:17:12+02', NULL, 1, 1, '2026-01-22 22:17:12+02', NULL, 'completed'),
(12, '2026-01-22 22:17:12+02', '2026-01-22 22:17:12+02', NULL, 1, 1, '2026-01-22 22:17:12+02', NULL, 'completed'),
(13, '2026-01-22 22:26:16+02', '2026-01-22 22:26:16+02', NULL, 1, 1, '2026-01-22 22:26:16+02', NULL, 'completed'),
(14, '2026-01-22 22:27:31+02', '2026-01-22 22:27:31+02', NULL, 1, 1, '2026-01-22 22:27:31+02', NULL, 'completed'),
(15, '2026-01-22 22:27:53+02', '2026-01-22 22:27:53+02', NULL, 1, 1, '2026-01-22 22:27:53+02', NULL, 'completed'),
(16, '2026-01-22 22:28:28+02', '2026-01-22 22:28:28+02', NULL, 1, 2, '2026-01-22 22:28:28+02', NULL, 'completed'),
(17, '2026-01-22 22:28:45+02', '2026-01-22 22:28:45+02', NULL, 1, 2, '2026-01-22 22:28:45+02', NULL, 'completed'),
(18, '2026-01-22 22:50:17+02', '2026-01-22 22:50:17+02', NULL, 1, 1, '2026-01-22 22:50:17+02', NULL, 'completed'),
(19, '2026-01-22 22:51:55+02', '2026-01-22 22:51:55+02', NULL, 1, 1, '2026-01-22 22:51:55+02', NULL, 'completed'),
(20, '2026-01-22 22:52:49+02', '2026-01-22 22:52:49+02', NULL, 1, 1, '2026-01-22 22:52:49+02', NULL, 'completed'),
(21, '2026-01-22 22:58:20+02', '2026-01-22 22:58:39+02', NULL, 1, 1, '2026-01-22 22:58:20+02', '2026-01-22 22:58:39+02', 'completed'),
(22, '2026-01-22 23:19:51+02', '2026-01-22 23:20:09+02', NULL, 1, 1, '2026-01-22 23:19:51+02', '2026-01-22 23:20:09+02', 'completed'),
(23, '2026-01-22 23:33:54+02', '2026-01-22 23:33:54+02', NULL, 1, 1, '2026-01-22 23:33:54+02', NULL, 'completed'),
(24, '2026-01-22 23:33:54+02', '2026-01-22 23:34:27+02', NULL, 1, 1, '2026-01-22 23:33:54+02', '2026-01-22 23:34:27+02', 'completed'),
(25, '2026-02-03 13:50:48+02', '2026-02-03 13:50:51+02', NULL, 1, 2, '2026-02-03 13:50:48+02', '2026-02-03 13:50:51+02', 'completed'),
(26, '2026-02-03 17:09:42+02', '2026-02-03 17:11:52+02', NULL, 1, 2, '2026-02-03 17:09:42+02', '2026-02-03 17:11:52+02', 'completed'),
(27, '2026-02-03 19:03:06+02', '2026-02-03 19:03:42+02', NULL, 3, 1, '2026-02-03 19:03:06+02', '2026-02-03 19:03:42+02', 'completed'),
(28, '2026-02-03 19:04:00+02', '2026-02-03 19:04:38+02', NULL, 1, 1, '2026-02-03 19:04:00+02', '2026-02-03 19:04:38+02', 'completed'),
(29, '2026-02-03 19:46:53+02', '2026-02-03 19:51:32+02', NULL, 1, 1, '2026-02-03 19:46:53+02', '2026-02-03 19:51:32+02', 'completed'),
(30, '2026-02-03 19:53:25+02', '2026-02-03 19:53:42+02', NULL, 3, 1, '2026-02-03 19:53:25+02', '2026-02-03 19:53:42+02', 'completed');

-- Conversations
INSERT INTO "conversations" ("id", "created_at", "updated_at", "deleted_at", "user1_id", "user2_id", "last_message_at") VALUES 
(1, '2026-02-04 11:08:28+02', '2026-02-04 11:08:28+02', NULL, 3, 3, '2026-02-04 11:08:28+02'),
(2, '2026-02-04 11:09:27+02', '2026-02-04 11:41:43+02', NULL, 1, 3, '2026-02-04 11:41:43+02');

-- Messages
INSERT INTO "messages" ("id", "created_at", "updated_at", "deleted_at", "conversation_id", "sender_id", "content", "is_read") VALUES 
(1, '2026-02-04 11:09:39+02', '2026-02-04 11:09:39+02', NULL, 2, 1, 'Hi Testov!', false),
(2, '2026-02-04 11:10:11+02', '2026-02-04 11:10:11+02', NULL, 2, 1, 'Wow', false),
(3, '2026-02-04 11:39:39+02', '2026-02-04 11:39:39+02', NULL, 2, 3, 'Hi', false),
(4, '2026-02-04 11:39:50+02', '2026-02-04 11:39:50+02', NULL, 2, 1, 'how are you', false),
(5, '2026-02-04 11:39:58+02', '2026-02-04 11:39:58+02', NULL, 2, 3, 'fine thanks', false),
(6, '2026-02-04 11:40:27+02', '2026-02-04 11:40:27+02', NULL, 2, 3, 'Training tomorrow?', false),
(7, '2026-02-04 11:41:43+02', '2026-02-04 11:41:43+02', NULL, 2, 1, 'I do not know if I will be able to tomorrow.', false);

-- Classes
INSERT INTO "classes" ("id", "name", "description", "trainer_id", "start_time", "duration", "max_capacity", "created_at", "difficulty_level") VALUES 
(7, 'Power Hour HIIT', 'A high-intensity interval training session designed to push your limits and burn maximum calories.', 1, '2026-01-19 15:40:56+02', 45, 15, '2026-01-18 15:40:56+02', NULL),
(8, 'Zen Yoga Flow', 'Focus on flexibility, balance, and mindful breathing in this restorative morning flow.', 1, '2026-01-20 15:40:56+02', 60, 25, '2026-01-18 15:40:56+02', NULL),
(9, 'Heavy Lifting 101', 'Master the fundamentals of the deadlift, squat, and bench press with expert form coaching.', 1, '2026-01-21 15:40:56+02', 90, 10, '2026-01-18 15:40:56+02', NULL),
(10, 'Power Hour HIIT', 'A high-intensity interval training session designed to push your limits and burn maximum calories.', 1, '2026-01-19 15:49:48+02', 45, 15, '2026-01-18 15:49:48+02', NULL),
(11, 'Zen Yoga Flow', 'Focus on flexibility, balance, and mindful breathing in this restorative morning flow.', 1, '2026-01-20 15:49:48+02', 60, 25, '2026-01-18 15:49:48+02', NULL),
(12, 'Heavy Lifting 101', 'Master the fundamentals of the deadlift, squat, and bench press with expert form coaching.', 1, '2026-01-21 15:49:48+02', 90, 10, '2026-01-18 15:49:48+02', NULL);

-- Enrollments
INSERT INTO "enrollments" ("id", "user_id", "class_id", "status", "created_at") VALUES 
(23, 1, 11, 'active', '2026-01-19 00:55:15+02'),
(24, 1, 8, 'active', '2026-01-19 00:57:08+02'),
(25, 1, 12, 'active', '2026-01-21 20:27:22+02'),
(26, 1, 9, 'active', '2026-01-21 20:27:27+02'),
(30, 3, 11, 'active', '2026-02-03 19:27:13+02'),
(32, 2, 10, 'active', '2026-02-03 19:27:54+02');

-- 3. CRITICAL: Reset the Auto-Increment IDs
-- Since we inserted IDs manually, we must tell Postgres to start counting from the highest ID.
SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 1)) FROM "users";
SELECT setval(pg_get_serial_sequence('exercises', 'id'), COALESCE(MAX(id), 1)) FROM "exercises";
SELECT setval(pg_get_serial_sequence('workout_plans', 'id'), COALESCE(MAX(id), 1)) FROM "workout_plans";
SELECT setval(pg_get_serial_sequence('workout_exercises', 'id'), COALESCE(MAX(id), 1)) FROM "workout_exercises";
SELECT setval(pg_get_serial_sequence('workout_sessions', 'id'), COALESCE(MAX(id), 1)) FROM "workout_sessions";
SELECT setval(pg_get_serial_sequence('conversations', 'id'), COALESCE(MAX(id), 1)) FROM "conversations";
SELECT setval(pg_get_serial_sequence('messages', 'id'), COALESCE(MAX(id), 1)) FROM "messages";
SELECT setval(pg_get_serial_sequence('classes', 'id'), COALESCE(MAX(id), 1)) FROM "classes";
SELECT setval(pg_get_serial_sequence('enrollments', 'id'), COALESCE(MAX(id), 1)) FROM "enrollments";