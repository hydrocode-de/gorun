export interface CitationFile {
    // Basic metadata
    Title: string;
    Version?: string;
    Message?: string;
    CffVersion?: string;
    CffType?: string;
    Abstract?: string;
    Keywords?: string[];

    // Repository information
    Repository?: Record<string, any>;
    RepositoryArtifact?: Record<string, any>;
    RepositoryCode?: {
        Scheme: string;
        Host: string;
        Path: string;
        ForceQuery: boolean;
        Fragment: string;
        OmitHost: boolean;
        Opaque: string;
        RawFragment: string;
        RawPath: string;
        RawQuery: string;
        User: null;
    };

    // License information
    License?: {
        Data: string[];
    };
    LicenseUrl?: Record<string, any>;

    // DOI information
    Doi?: {
        General: string;
        DirectoryIndicator: string;
        RegistrantCode: string;
    };

    // URL information
    Url?: {
        Scheme: string;
        Host: string;
        Path: string;
        ForceQuery: boolean;
        Fragment: string;
        OmitHost: boolean;
        Opaque: string;
        RawFragment: string;
        RawPath: string;
        RawQuery: string;
        User: null;
    };

    // Authors/Contributors
    Authors?: Author[];
    Contact?: Author[];

    // Dates
    DateReleased?: Record<string, any>;
    Commit?: string;

    // References
    References?: Reference[];
    Identifiers?: null;

    // Preferred citation
    PreferredCitation?: {
        Title?: string;
        Authors?: Author[] | null;
        Contact?: Author[] | null;
        Abstract?: string;
        Keywords?: string[] | null;
        Version?: string;
        Year?: number;
        Month?: number;
        Volume?: number;
        Issue?: string;
        Pages?: string;
        Journal?: string;
        Publisher?: Entity;
        Institution?: Entity;
        Location?: Entity;
        Conference?: Entity;
        DatabaseProvider?: Entity;
        Doi?: {
            General: string;
            DirectoryIndicator: string;
            RegistrantCode: string;
        };
        CollectionDoi?: {
            General: string;
            DirectoryIndicator: string;
            RegistrantCode: string;
        };
        CollectionTitle?: string;
        CollectionType?: string;
        Copyright?: string;
        DataType?: string;
        Database?: string;
        Department?: string;
        Edition?: string;
        Editors?: Author[] | null;
        EditorsSeries?: Author[] | null;
        End?: string;
        Entry?: string;
        Filename?: string;
        Format?: string;
        Identifiers?: null;
        Isbn?: string;
        Issn?: string;
        IssueDate?: string;
        IssueTitle?: string;
        Languages?: string[] | null;
        License?: {
            Data: string[] | null;
        };
        LicenseUrl?: Record<string, any>;
        LocEnd?: string;
        LocStart?: string;
        Medium?: string;
        Nihmsid?: string;
        Notes?: string;
        Number?: string;
        NumberVolumes?: string;
        PatentStates?: string[] | null;
        Pmcid?: string;
        Recipients?: Author[] | null;
        Repository?: Record<string, any>;
        RepositoryArtifact?: Record<string, any>;
        RepositoryCode?: Record<string, any>;
        Scope?: string;
        Section?: string;
        Senders?: Author[] | null;
        Start?: string;
        Status?: string;
        Term?: string;
        ThesisType?: string;
        Translators?: Author[] | null;
        ReferenceType?: string;
        Url?: Record<string, any>;
        VolumeTitle?: string;
        YearOriginal?: number;
        DateAccessed?: Record<string, any>;
        DateDownloaded?: Record<string, any>;
        DatePublished?: Record<string, any>;
        Abbreviation?: string;
    };
}

export interface Author {
    IsEntity: boolean;
    IsPerson: boolean;
    Entity?: Entity;
    Person?: Person;
}

export interface Entity {
    Name: string;
    Address?: string;
    City?: string;
    Country?: string;
    Email?: string;
    Fax?: string;
    Location?: string;
    Orcid?: string;
    PostCode?: string;
    Region?: string;
    Tel?: string;
    Website?: Record<string, any>;
    Alias?: string;
    DateStart?: Record<string, any>;
    DateEnd?: Record<string, any>;
}

export interface Person {
    Family: string;
    GivenNames?: string;
    Address?: string;
    Affiliation?: string;
    City?: string;
    Country?: string;
    Email?: string;
    Fax?: string;
    NameParticle?: string;
    NameSuffix?: string;
    Orcid?: string;
    PostCode?: string;
    Region?: string;
    Tel?: string;
    Website?: Record<string, any>;
    Alias?: string;
}

export interface Reference {
    Type: string;
    Title: string;
    Authors?: Author[];
    DOI?: string;
    URL?: string;
    Year?: number;
    Month?: number;
    Day?: number;
    Journal?: string;
    Volume?: string;
    Issue?: string;
    FirstPage?: string;
    LastPage?: string;
    Publisher?: string;
    Place?: string;
    Status?: string;
}