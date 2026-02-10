# WebSocket Chat with Private Messages 💬

## المميزات ✨

1. **Public Messages** - رسائل عامة لجميع المستخدمين
2. **Private Messages** - رسائل خاصة بين مستخدمين
3. **Login System** - كل مستخدم عندو username
4. **Real-time Communication** - تواصل فوري عبر WebSocket

## البنية 🏗️

```
ws-private/
├── main.go          # نقطة الدخول
├── manager.go       # إدارة الاتصالات والـ events
├── client.go        # معالجة كل client
├── event.go         # تعريف الأحداث
├── go.mod           # Go modules
├── go.sum           
└── frontend/
    └── index.html   # واجهة المستخدم
```

## كيفاش يخدم 🚀

### 1. Login Flow
```
Client → send "login" event → Server
Server → يخزن username في map[string]*Client
```

### 2. Public Message
```
Client A → "send_message" → Server
Server → broadcast لجميع الـ clients
```

### 3. Private Message
```
Client A → "private_message" {to: "userB"} → Server
Server → يجيب Client B من map
Server → يسيفط الرسالة لـ Client B فقط
```

## التشغيل 🏃

```bash
# نزل الـ dependencies
go mod download

# شغل الـ server
go run .

# دخل للموقع
http://localhost:8080
```

## الاستعمال 📝

1. **Login**: دخل username ديالك
2. **Public Message**: اكتب رسالة و كبس "Send to All"
3. **Private Message**: حدد المستقبل و اكتب الرسالة

## Event Types 📡

### Login Event
```json
{
  "type": "login",
  "payload": {
    "username": "ahmed"
  }
}
```

### Public Message Event
```json
{
  "type": "send_message",
  "payload": {
    "message": "Hello everyone!",
    "from": "ahmed"
  }
}
```

### Private Message Event
```json
{
  "type": "private_message",
  "payload": {
    "message": "Hey, how are you?",
    "from": "ahmed",
    "to": "sara"
  }
}
```

### New Message Event (من Server)
```json
{
  "type": "new_message",
  "payload": {
    "message": "Hello!",
    "from": "ahmed",
    "sent": "14:30:25"
  }
}
```

## الأمان 🔒

⚠️ هاد الكود للتعلم فقط! في production خاصك:

- ✅ Authentication صحيح (JWT, sessions)
- ✅ Validation ديال usernames (unique, sanitized)
- ✅ Rate limiting
- ✅ Message encryption (TLS/SSL)
- ✅ Input sanitization (XSS protection)
- ✅ CORS configuration صحيح

## التحسينات الممكنة 🎯

- [ ] قائمة بالمستخدمين المتصلين
- [ ] History ديال الرسائل
- [ ] Typing indicators
- [ ] Read receipts
- [ ] File sharing
- [ ] Group chats
- [ ] Message reactions
- [ ] User status (online/offline)

## الهيكل التقني 🔧

### Manager
- `clients map[*Client]bool` - جميع الاتصالات
- `users map[string]*Client` - username → Client
- Thread-safe باستعمال `sync.RWMutex`

### Client
- `username string` - اسم المستخدم
- `egress chan Event` - قناة لإرسال الرسائل
- `ReadMessage()` - goroutine لقراءة الرسائل
- `writeMessages()` - goroutine لإرسال الرسائل

### Event System
- Event-based architecture
- Handler pattern لمعالجة الأحداث
- JSON serialization

## المساهمة 🤝

أي اقتراحات أو تحسينات؟ مرحبا بيك!

---
Made with ❤️ using Go & WebSocket
