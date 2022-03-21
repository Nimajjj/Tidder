CREATE TABLE Accounts(
   id_account INT AUTO_INCREMENT NOT NULL,
   name VARCHAR(25) NOT NULL,
   email VARCHAR(50) NOT NULL,
   hashed_password VARCHAR(25) NOT NULL,
   birth_date DATE NOT NULL,
   creation_date DATE NOT NULL,
   karma VARCHAR(50) NOT NULL DEFAULT(0),
   profile_picture VARCHAR(255) NOT NULL DEFAULT("/images/profiles/default.png"),
   PRIMARY KEY(id_account),
   UNIQUE(name),
   UNIQUE(email)
);

CREATE TABLE Subjects(
   id_subject INT AUTO_INCREMENT NOT NULL,
   name VARCHAR(25) NOT NULL,
   profile_picture VARCHAR(255) NOT NULL,
   id_owner INT NOT NULL,
   nsfw BOOLEAN NOT NULL DEFAULT(FALSE),
   PRIMARY KEY(id_subject),
   UNIQUE(name),
   FOREIGN KEY(id_owner) REFERENCES Accounts(id_account)
);

CREATE TABLE Posts(
   id_post INT AUTO_INCREMENT NOT NULL,
   title VARCHAR(255) NOT NULL,
   media_url VARCHAR(255),
   content TEXT NOT NULL,
   creation_date DATETIME NOT NULL,
   upvotes INT NOT NULL DEFAULT(0),
   downvotes INT NOT NULL DEFAULT(0),
   nsfw BOOLEAN NOT NULL DEFAULT(FALSE),
   redacted BOOLEAN NOT NULL DEFAULT(FALSE),
   pinned BOOLEAN NOT NULL DEFAULT(FALSE),
   id_subject INT NOT NULL,
   id_author INT NOT NULL,
   PRIMARY KEY(id_post),
   FOREIGN KEY(id_subject) REFERENCES Subjects(id_subject),
   FOREIGN KEY(id_author) REFERENCES Accounts(id_account)
);

CREATE TABLE Comments(
   id_comment INT AUTO_INCREMENT NOT NULL,
   content TEXT NOT NULL,
   creation_date DATETIME NOT NULL,
   upvotes INT NOT NULL DEFAULT(0),
   downvotes INT NOT NULL DEFAULT(0),
   redacted BOOLEAN NOT NULL DEFAULT(FALSE),
   id_author INT NOT NULL,
   response_to_id INT NOT NULL DEFAULT(-1),
   id_post INT NOT NULL,
   PRIMARY KEY(id_comment),
   FOREIGN KEY(id_author) REFERENCES Accounts(id_account),
   FOREIGN KEY(response_to_id) REFERENCES Comments(id_comment),
   FOREIGN KEY(id_post) REFERENCES Posts(id_post)
);

CREATE TABLE Subject_Access(
   id_subject_access INT AUTO_INCREMENT NOT NULL,
   pin_post BOOLEAN NOT NULL DEFAULT(FALSE),
   remove_post BOOLEAN NOT NULL DEFAULT(FALSE),
   ban_user BOOLEAN NOT NULL DEFAULT(FALSE),
   create_role BOOLEAN NOT NULL DEFAULT(FALSE),
   give_role BOOLEAN NOT NULL DEFAULT(FALSE),
   delete_role BOOLEAN NOT NULL DEFAULT(FALSE),
   PRIMARY KEY(id_subject_access)
);

CREATE TABLE Global_Roles(
   id_global_role INT AUTO_INCREMENT NOT NULL,
   name VARCHAR(25) NOT NULL,
   PRIMARY KEY(id_global_role),
   UNIQUE(name)
);

CREATE TABLE Global_Access(
   id_global_access INT AUTO_INCREMENT NOT NULL,
   name VARCHAR(25) NOT NULL,
   PRIMARY KEY(id_global_access)
);

CREATE TABLE Subject_Roles(
   id_subject_role INT AUTO_INCREMENT NOT NULL,
   name VARCHAR(25) NOT NULL,
   id_subject INT NOT NULL,
   id_subject_access INT NOT NULL,
   PRIMARY KEY(id_subject_role),
   UNIQUE(name),
   FOREIGN KEY(id_subject) REFERENCES Subjects(id_subject),
   FOREIGN KEY(id_subject_access) REFERENCES Subject_Access(id_subject_access)
);

CREATE TABLE Subscribe_to_subject(
   id_account INT NOT NULL,
   id_subject INT NOT NULL,
   PRIMARY KEY(id_account, id_subject),
   FOREIGN KEY(id_account) REFERENCES Accounts(id_account),
   FOREIGN KEY(id_subject) REFERENCES Subjects(id_subject)
);

CREATE TABLE Global_Roles_Management(
   id_account INT NOT NULL,
   id_global_role INT NOT NULL,
   id_global_access INT NOT NULL,
   PRIMARY KEY(id_account, id_global_role, id_global_access),
   FOREIGN KEY(id_account) REFERENCES Accounts(id_account),
   FOREIGN KEY(id_global_role) REFERENCES Global_Roles(id_global_role),
   FOREIGN KEY(id_global_access) REFERENCES Global_Access(id_global_access)
);

CREATE TABLE Is_Ban(
   id_account INT NOT NULL,
   id_subject INT NOT NULL,
   PRIMARY KEY(id_account, id_subject),
   FOREIGN KEY(id_account) REFERENCES Accounts(id_account),
   FOREIGN KEY(id_subject) REFERENCES Subjects(id_subject)
);

CREATE TABLE Has_Subject_Role(
   id_account INT NOT NULL,
   id_subject INT NOT NULL,
   id_subject_role INT NOT NULL,
   PRIMARY KEY(id_account, id_subject, id_subject_role),
   FOREIGN KEY(id_account) REFERENCES Accounts(id_account),
   FOREIGN KEY(id_subject) REFERENCES Subjects(id_subject),
   FOREIGN KEY(id_subject_role) REFERENCES Subject_Roles(id_subject_role)
);
