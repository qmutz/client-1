//
//  KBEnvironment.h
//  Keybase
//
//  Created by Gabriel on 4/22/15.
//  Copyright (c) 2015 Gabriel Handford. All rights reserved.
//

#import <Foundation/Foundation.h>

typedef NS_ENUM (NSInteger, KBEnv) {
  KBEnvManual = 1,
  KBEnvLocalhost,
  KBEnvKeybaseIO,
};

@interface KBEnvironment : NSObject

@property (readonly) NSString *homeDir;
@property (readonly) NSString *host;
@property (readonly, getter=isDebugEnabled) BOOL debugEnabled;
@property (readonly) NSString *mountDir;
@property (readonly) NSString *sockFile;
@property (readonly) NSString *identifier;
@property (readonly) NSString *launchdLabelService;
@property (readonly) NSString *launchdLabelKBFS;
@property (readonly) NSString *title;
@property (readonly) NSString *info;
@property (readonly) NSImage *image;
@property (readonly) BOOL canRunFromXCode;

- (instancetype)initWithEnv:(KBEnv)env;

+ (instancetype)env:(KBEnv)env;

- (NSDictionary *)launchdPlistDictionaryForService;
- (NSDictionary *)launchdPlistDictionaryForKBFS;

@end
