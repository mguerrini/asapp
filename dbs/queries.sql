//Get Messages
SELECT m.*,
       CASE
           WHEN txt.Id IS NOT NULL THEN 'text'
           WHEN vd.Id IS NOT NULL THEN 'video'
           WHEN img.Id IS NOT NULL THEN 'image'
       END as ContentType,
       txt.Id AS text_Id,
       txt.Text AS text_Text,
       vd.Id AS video_Id,
       vd.Url AS video_Url,
       vd.Source AS video_Source,
       img.Id AS image_Id,
       img.Url AS image_Url,
       img.Height AS image_Height,
       img.Width AS image_Width

FROM Messages m
         LEFT JOIN TextContent txt ON m.Id = txt.MessageId
         LEFT JOIN VideoContent vd ON m.Id = vd.MessageId
         LEFT JOIN ImageContent img ON m.Id = img.MessageId

WHERE @startId <= m.Id AND m.Id < @finishId AND @recipient = m.RecipientId

//Insert Message
INSERT INTO Messages (
    SenderId,
    RecipientId,
    Timestamp
)
VALUES (
    @senderId,
    @recipientId,
    datetime('now')
    );


// IMAGE
INSERT INTO ImageContent (
    MessageId,
    Url,
    Height,
    Width
)
VALUES (
    @MessageId,
    @Url,
    @Height,
    @Width
);



// TEXT
INSERT INTO TextContent (
    MessageId,
    Text
)
VALUES (
    @MessageId,
    @Text
);



//VIDEO

INSERT INTO VideoContent (
     MessageId,
     Url,
     Source
 )
 VALUES (
     @MessageId,
     @Url,
     @Source
 );


// User
INSERT INTO Users (
    Username,
    Password
)
VALUES (
   @Username,
   @Password
);


// Get Users
SELECT Id,
       Username,
       Password
FROM Users
WHERE (@userId is null OR @userId = Id) AND (@userName is null OR @userName = Username)
