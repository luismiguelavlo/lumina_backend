# Lumina Library — Landing page draft (what students can do)

Student-facing copy only: **things you can open, view, and track**—not staff workflows. Based on `functionalities/README.md` (public ranking, catalog, student profile, loans as something you *see*, reputation/badges as feedback you *get*).

---

## 1. Hero — You’re here to read

**Suggested headline**  
See the collection. See the leaders. See your progress.

**Subheadline**  
Lumina is your campus library online: browse what’s on the shelves, check who’s topping the reader board, and keep your loans and stats in one place—built around students, not back-office screens.

**Microcopy (optional under CTAs)**  
Top 3 readers · Book catalog · Your profile & history

**Primary CTA**  
Explore the book catalog

**Secondary CTA**  
View the Top 3

---

## 2. Meet the Top 3 — Campus readers in the spotlight

**Section headline**  
See who’s leading the pack.

**Body**  
Anyone can open the **Top 3** and see the three highest-ranked readers by reputation points—names and scores, no login required. It’s a quick, fun snapshot of who’s been showing up for the library habit lately, and a nudge to see where you might land if you keep returning on time.

**What you can do**  
- Open the Top 3 anytime (public).  
- Share it with friends or your study group as a light challenge.  
- Use it as a conversation starter—“who’s on the board this week?”

---

## 3. Explore the book catalog — What’s in the library

**Section headline**  
Browse books like you browse a playlist.

**Body**  
Search and flip through the **book catalog**: covers, titles, authors, genres, and the story behind each title—synopsis, year, length, and where to find it in the building. Check **availability** so you know before you go whether a copy is free or already out with someone else.

**Coming soon**  
We’re **upgrading the way you explore the catalog**—clearer discovery, smoother browsing, and a better feel for what’s new and what fits your courses. The collection is the same; the student experience is getting better.

**What you can do**  
- Search by title, author, or ISBN.  
- Scan genres and authors when you want a recommendation path.  
- See shelf location and live availability.  
- Plan your visit with a short list of titles you actually want.

---

## 4. Your place on the board — Full leaderboard

**Section headline**  
Find your rank—not just the podium.

**Body**  
Beyond the Top 3, you can open the **global leaderboard** and scroll other positions. Plug in **one** identifier the library already knows—your student code, institutional email, or internal student id—and the app can highlight **your** row: rank, points, and reading streak signals when they’re available. No mystery scores: you see how you compare and what to improve.

**What you can do**  
- Browse ranks past the first three spots.  
- Look up *your* standing with a single id field (as documented for the product).  
- Pair it with your profile to connect points with real loans and returns.

---

## 5. Your profile — Loans and history you can actually see

**Section headline**  
Everything you’re reading, in one student view.

**Body**  
Your **profile** is the home for your library life: avatar, program, and member-since vibes up top, then the numbers that matter—books finished, active loans, streaks, and reputation points. Below that, **loan history** reads like a timeline: what’s **out now** (“currently reading”) and what you’ve **already brought back**, with covers and authors so you never confuse two similar titles.

**What you can do**  
- See active vs returned loans at a glance.  
- Track deadlines and patterns without digging through email.  
- Treat it as your personal reading resume on campus.

---

## 6. Points & streaks — Feedback for showing up

**Section headline**  
Good habits show up in your score.

**Body**  
Returns feed your reputation: **on-time** behavior adds points and keeps your **streak** alive; late returns are recorded in a straightforward way so the rules are the same for everyone. You’re not managing a spreadsheet—the app reflects what you already did at the library.

**What you can do**  
- Watch points and streaks update as you return books.  
- Use streaks as a gentle nudge, not a lecture.  
- Understand that the score is shared context for the leaderboard—not a hidden staff-only report.

---

## 7. Badges — Goals you can collect

**Section headline**  
Unlock badges as you grow as a reader.

**Body**  
The **badge gallery** on your profile shows what exists, what you’ve earned, and what’s still locked—with icons and short criteria so progress feels tangible. It’s optional motivation: some students ignore it, others collect every one.

**What you can do**  
- See earned vs locked badges in one grid.  
- Read criteria and know what “complete” looks like.  
- Show off progress to friends if that’s your thing.

---

## 8. The library, for you — Why Lumina exists

**Section headline**  
A campus library that feels alive.

**Body**  
Lumina isn’t a pile of forms—it’s the digital side of the stacks: **discover books**, **see who’s reading with you** (via public rankings), and **carry your history** in a profile that respects the student experience. Staff keep the shelves accurate; you get the view that helps you study and explore.

**Final CTA**  
Open the catalog, peek at the Top 3, and find your name on the board.

---

## 9. Footer — Quick links (student lens)

**Explore**  
Book catalog · Top 3 · Leaderboard · Your profile (when you have access)

**Library**  
Hours · Visit the building · Help / contact

**Brief legal line**  
Data use according to your institution’s policies.

---

## Internal notes (do not publish on the landing)

- `GET /api/ranking/top3` and `GET /api/ranking/leaderboard` are **public** in spec—ideal hero/embedded widgets for students without account flows.  
- Catalog UX is flagged as **changing soon** in this copy; adjust the “Coming soon” blurb when the new browse ships.  
- Profile and full catalog depth may still depend on how the frontend wires auth; keep student promises aligned with what’s actually exposed in the app.
