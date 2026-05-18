# Product Overview

SWU OSR (Open Source Repository) is a campus open-source repository platform for SWU (Sains dan Teknologi Walisongo) university in Purwokerto, Indonesia. It is an HMPSTI student organization project.

## What It Does

- Links student academic identities from SIAKAD (campus academic system) with their GitHub accounts
- Aggregates GitHub activity via webhooks and displays it on a public dashboard
- Enforces pseudonym-first culture: students interact under aliases, only faculty/admin can see real identities
- Provides an internal discussion forum for campus-scoped collaboration

## Current State (SWAAP)

The existing codebase is a "Smart Wrapper Academic Portal" — a Go API that proxies authentication and data extraction from the legacy PHP-based SIAKAD system at `smartone.smart-service.co.id`. It simulates browser sessions to log students in, fetch schedules, attendance records, and submit attendance. A Flutter web/mobile frontend consumes this API.

## Planned State (SWU OSR Platform)

The platform evolves through 5 phases:
1. Authentication & Identity Binding (SIAKAD proxy + GitHub OAuth)
2. Onboarding & Showcase Selection (students curate public repos with academic tags)
3. Aggregator Engine (GitHub webhooks capture activity in real-time)
4. Public Dashboard & Identity Validation (role-based access to real identities)
5. Internal Discussion Forum (campus-scoped threads and comments)

## Key Domain Concepts

- **NIM**: Student identification number (Nomor Induk Mahasiswa)
- **Alias**: Pseudonym chosen by student for public display
- **Showcase Repo**: GitHub repository selected by student for platform display
- **Academic Tag**: Classification label (coursework, thesis, hackathon, personal_research, team_project)
- **Activity Log**: Record of a GitHub event (push, pull_request, release) captured via webhook
- **SIAKAD**: The campus academic information system (external, PHP-based)
