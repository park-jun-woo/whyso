## go
[billing]HoldEscrow,ReleaseFunds

## ssac
[service/auth]Login,Register
[service/gig]ApproveWork,CreateGig,GetGig,ListGigs,PublishGig,RaiseDispute,SubmitProposal,SubmitWork
[service/proposal]AcceptProposal,RejectProposal

## openapi
[api]AcceptProposal,ApproveWork,CreateGig,GetGig,ListGigs,Login,PublishGig,RaiseDispute,Register,RejectProposal,SubmitProposal,SubmitWork

## sql
[db/queries]GigCreate,GigFindByID,GigList,GigUpdateFreelancerID,GigUpdateStatus,ProposalCreate,ProposalFindByID,ProposalUpdateStatus,TransactionCreate,UserCreate,UserFindByEmail,UserFindByID

## rego
[policy]allow

## gherkin
[scenario]Happy Path - Full Gig Lifecycle,Invalid State - Cannot approve work when gig is in open state,Unauthorized Access - Freelancer B cannot submit work on Freelancer A gig

## stml
[frontend]data-action:ApproveWork,data-action:PublishGig,data-action:RaiseDispute,data-action:SubmitWork,data-fetch:GetGig,data-fetch:ListGigs

## mermaid
[states]accepted,completed,disputed,draft,in_progress,open,pending,rejected,under_review
