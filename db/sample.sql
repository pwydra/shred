insert into shred_user (
    user_uuid, first_name, last_name, email, created_by
) values ('f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1f', 'John', 'Doe', 'jdoe@gmail.com', 'f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1f');

insert into muscle_type (
    muscle_code, muscle_name, muscle_description, 
    muscle_group, created_by
) values ('QUAD', 'Quadriceps', 'The quadriceps are a group of muscles located at the front of the thigh.', 
    'Legs', 'f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1f');

insert into category_type (
    category_code, category_name, category_description, 
    created_by
) values ('STRENGTH', 'Strength Training', 'Exercises that improve strength and endurance.', 
    'f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1f');

insert into apparatus_type (
    apparatus_code, apparatus_name, apparatus_description, 
    created_by
) values ('BARBELL', 'Barbell', 'A long bar with weights on either end used for strength training.', 
    'f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1f');

insert into license (
    license_short_name, license_full_name, url, 
    created_by
) values ('CC_BY', 'Creative Commons Attribution', 'https://creativecommons.org/licenses/by/4.0/', 
    'f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1f');

insert into exercise (
    exercise_name, exercise_description, instructions, cues,
    video_url, category_code, license_short_name, license_author, 
    created_by
) values (
    'Squat', 'A lower body exercise that targets the quadriceps, hamstrings, and glutes.',
    'Stand with feet shoulder-width apart, lower your body as if sitting back into a chair, then return to standing.',
    'Keep your chest up, knees behind toes, and weight on your heels.'
    , 'https://example.com/squat_video.mp4', 'STRENGTH', 'CC_BY', 'John Doe',
    'f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1f'
);

commit;